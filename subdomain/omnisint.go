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
	return &Threatcrowd{
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

func (o Omnisint) Enumeration(domain string) (map[string]struct{}, error) {
	result := make(map[string]struct{})
	urlAddress := fmt.Sprintf(o.Url+"%s", domain)
	resp, err := http.Get(urlAddress)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	var subdomains []string
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&subdomains)
	for _, subdomainAddress := range subdomains {
		subdomain, err := url.Parse("https://" + subdomainAddress)
		if err != nil {
			continue
		}
		result[subdomain.Hostname()] = struct{}{}
	}

	return result, nil
}
