package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
)

type JLDC struct {
	Url string
}

type JLDCRespose struct {
	subdomains []string
}

func NewJLDC() SubdomainFinderInterface {
	return JLDC{Url: constants.JLDCUrl}
}

func (j JLDC) IsPaidProvider() bool {
	return false
}

func (j JLDC) SetAuth(token string) {}

func (j JLDC) Name() string {
	return "jldc"
}

func (j JLDC) Enumeration(domain string) (map[string]struct{}, error) {
	result := make(map[string]struct{})
	fetchURL := fmt.Sprintf(j.Url+"/%s", domain)
	resp, err := http.Get(fetchURL)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	var jldcRes JLDCRespose
	err = dec.Decode(&jldcRes)
	for _, r := range jldcRes.subdomains {
		subdomain, err := url.Parse("https://" + r)
		if err != nil {
			continue
		}
		result[subdomain.Hostname()] = struct{}{}
	}

	return result, nil
}
