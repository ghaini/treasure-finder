package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
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
		Url: constants.CertUrl,
	}
}

func (c Crt) Enumeration(domain string, subdomains chan<- string) {
	urlAddress := fmt.Sprintf(c.Url+"?q=%%25.%s&output=json", domain)
	resp, err := http.Get(urlAddress)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var crtResponse []crtResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err := json.Unmarshal(body, &crtResponse); err != nil {
		return
	}
	for _, crt := range crtResponse {
		names := strings.Fields(crt.Name)
		for _, name := range names {
			subdomain, err := url.Parse(name)
			if err != nil {
				continue
			}
			subdomains <- subdomain.Hostname()
		}
	}
}
