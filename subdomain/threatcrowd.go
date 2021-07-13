package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
)

type Threatcrowd struct {
	Url string
}

type threatcrowdResponse struct {
	Subdomains []string `json:"subdomains"`
}

func NewThreatcrowd() SubdomainFinderInterface {
	return &Threatcrowd{
		Url: constants.ThreatcrowdUrl,
	}
}

func (t Threatcrowd) IsPaidProvider() bool {
	return false
}

func (t Threatcrowd) Name() string {
	return "threatcrowd"
}

func (t Threatcrowd) SetAuth(token string) {
	return
}

func (t Threatcrowd) Enumeration(domain string) (map[string]struct{}, error) {
	result := make(map[string]struct{})
	urlAddress := fmt.Sprintf(t.Url+"/domain/report/?domain=%s", domain)
	resp, err := http.Get(urlAddress)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	var threatcrowdRes threatcrowdResponse
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&threatcrowdRes)
	for _, subdomainAddress := range threatcrowdRes.Subdomains {
		subdomain, err := url.Parse("https://" + subdomainAddress)
		if err != nil {
			continue
		}
		result[subdomain.Hostname()] = struct{}{}
	}

	return result, nil
}
