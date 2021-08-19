package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
)

type Omnisint struct {
	Url string
}

func NewOmnisint() SubdomainFinderInterface {
	return &Omnisint{
		Url: constants.OmnisintUrl,
	}
}

func (o Omnisint) IsPaidProvider() bool {
	return false
}

func (o Omnisint) Name() string {
	return "omnisint"
}

func (o Omnisint) SetAuth(token string) {
	return
}

func (o Omnisint) GetAuth() string {return ""}

func (o Omnisint) Enumeration(domain string) (result map[string]struct{}, statusCode int, err error) {
	urlAddress := fmt.Sprintf(o.Url+"%s", domain)
	resp, err := http.Get(urlAddress)
	if err != nil {
		return result, 500, err
	}
	defer resp.Body.Close()

	var subdomains []string
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&subdomains)
	if err != nil {
		return result, resp.StatusCode, err
	}
	for _, subdomainAddress := range subdomains {
		subdomain, err := url.Parse("https://" + subdomainAddress)
		if err != nil {
			continue
		}
		result[subdomain.Hostname()] = struct{}{}
	}

	return result, resp.StatusCode, nil
}
