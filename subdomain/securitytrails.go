package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
)

type Securitytrails struct {
	Url   string
	token string
}
type SecuritytrailsResponse struct {
	Subdomains []string `json:"subdomains"`
}

func NewSecuritytrails() SubdomainFinderInterface {
	return &Securitytrails{
		Url: constants.SecuritytrailsUrl,
	}
}

func (s *Securitytrails) IsPaidProvider() bool {
	return true
}

func (s *Securitytrails) SetAuth(token string) {
	s.token = token
}

func (s *Securitytrails) Name() string {
	return "securitytrails"
}

func (s *Securitytrails) Enumeration(domain string) (map[string]struct{}, error) {
	result := make(map[string]struct{})
	urlAddress := fmt.Sprintf(s.Url+"/domain/%s/subdomains?apikey=%s", domain, s.token)
	resp, err := http.Get(urlAddress)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	var securitytrailsRes SecuritytrailsResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&securitytrailsRes)

	for _, subdomainUrl := range securitytrailsRes.Subdomains {
		address := "https://" + subdomainUrl + "." + domain
		subdomain, err := url.Parse(address)
		if err != nil {
			continue
		}
		result[subdomain.Hostname()] = struct{}{}
	}

	return result, nil
}
