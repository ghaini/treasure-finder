package subdomain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Crt struct {
	Url string
}

type crtResponse struct {
	Name string `json:"name_value"`
}

func NewCrt() SubdomainFinderInterface {
	return &Crt{
		Url: "https://crt.sh/",
	}
}

func (c Crt) Enumeration(domain string) (map[string]struct{}, error) {

	result := make(map[string]struct{})
	urlAddress := fmt.Sprintf(c.Url+"?q=%%25.%s&output=json", domain)
	resp, err := http.Get(urlAddress)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	var crtResponse []crtResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal(body, &crtResponse); err != nil {
		return result, err
	}
	for _, crt := range crtResponse {
		names := strings.Fields(crt.Name)
		for _, name := range names {
			subdomain, err := url.Parse(name)
			if err != nil {
				continue
			}
			result[subdomain.Hostname()] = struct{}{}
		}
	}

	return result, err
}
