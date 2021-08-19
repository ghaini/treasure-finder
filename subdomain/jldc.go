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

func NewJLDC() SubdomainFinderInterface {
	return JLDC{Url: constants.JLDCUrl}
}

func (j JLDC) IsPaidProvider() bool {
	return false
}

func (j JLDC) SetAuth(token string) {}

func (j JLDC) GetAuth() string {return ""}

func (j JLDC) Name() string {
	return "jldc"
}

func (j JLDC) Enumeration(domain string) (result map[string]struct{}, statusCode int, err error) {
	fetchURL := fmt.Sprintf(j.Url+"/%s", domain)
	resp, err := http.Get(fetchURL)
	if err != nil {
		return result, 500, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var jldcRes []string
	err = dec.Decode(&jldcRes)
	for _, r := range jldcRes {
		subdomain, err := url.Parse("https://" + r)
		if err != nil {
			continue
		}
		result[subdomain.Hostname()] = struct{}{}
	}

	return result, resp.StatusCode, nil
}
