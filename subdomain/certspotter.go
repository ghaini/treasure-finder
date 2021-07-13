package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
)

type Certspotter struct {
	Url string
}

type certspotterResponse struct {
	DNSNames []string `json:"dns_names"`
}

func NewCertspotter() SubdomainFinderInterface {
	return &Certspotter{
		Url: constants.CertspotterUrl,
	}
}

func (c Certspotter) IsPaidProvider() bool {
	return false
}

func (c Certspotter) SetAuth(token string) {
	return
}

func (c Certspotter) Name() string {
	return "certpotter"
}


func (c Certspotter) Enumeration(domain string) (map[string]struct{}, error) {
	result := make(map[string]struct{})
	urlAddress := fmt.Sprintf(c.Url+"/certs?domain=%s", domain)
	resp, err := http.Get(urlAddress)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	var certspotterResp []certspotterResponse
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&certspotterResp)
	for _, csr := range certspotterResp {
		for _, address := range csr.DNSNames {
			subdomain, err := url.Parse(address)
			if err != nil {
				continue
			}
			result[subdomain.Hostname()] = struct{}{}
		}
	}

	return result, nil
}