package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
)

type Alienvault struct {
	Url string
}

type alienvaultResponse struct {
	PassiveDNS []struct {
		Hostname string `json:"hostname"`
	} `json:"passive_dns"`
}

func NewAlienvault() SubdomainFinderInterface {
	return Alienvault{Url: constants.AlienvaultUrl}
}

func (a Alienvault) IsPaidProvider() bool {
	return false
}

func (a Alienvault) SetAuth(token string) {}

func (a Alienvault) GetAuth() string {return ""}

func (a Alienvault) Name() string {
	return "alienvault"
}

func (a Alienvault) Enumeration(domain string) (result map[string]struct{}, statusCode int, err error) {
	fetchURL := fmt.Sprintf(a.Url+"/%s/passive_dns", domain)
	resp, err := http.Get(fetchURL)
	if err != nil {
		return result, 500, err
	}

	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var alienvaultRes alienvaultResponse
	err = dec.Decode(&alienvaultRes)
	for _, r := range alienvaultRes.PassiveDNS {
		subdomain, err := url.Parse("https://" + r.Hostname)
		if err != nil {
			continue
		}
		result[subdomain.Hostname()] = struct{}{}
	}

	return result, resp.StatusCode,nil
}
