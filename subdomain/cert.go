package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Crt struct {
	Url string
}

type crtResponse struct {
	Name string `json:"name_value"`
}

func NewCrt() SubdomainFinderInterface {
	return &Crt{
		Url: constants.CertUrl,
	}
}

func (c Crt) IsPaidProvider() bool {
	return false
}

func (c Crt) Name() string {
	return "crt"
}

func (c Crt) SetAuth(token string) {
	return
}

func (c Crt) GetAuth() string {return ""}

func (c Crt) Enumeration(domain string) (result map[string]struct{}, statusCode int, err error){
	result = make(map[string]struct{})
	urlAddress := fmt.Sprintf(c.Url+"?q=%%25.%s&output=json", domain)
	client := http.Client{
		Timeout: 5 * time.Minute,
	}
	resp, err := client.Get(urlAddress)
	if err != nil {
		return result, 500, err
	}
	defer resp.Body.Close()

	var crtResponse []crtResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, resp.StatusCode, err
	}

	if err = json.Unmarshal(body, &crtResponse); err != nil {
		return result, resp.StatusCode, err
	}
	for _, crt := range crtResponse {
		names := strings.Fields(crt.Name)
		for _, name := range names {
			subdomain, err := url.Parse("https://" + name)
			if err != nil {
				continue
			}
			result[subdomain.Hostname()] = struct{}{}
		}
	}

	return result, resp.StatusCode, err
}
