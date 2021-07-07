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

func (c Certspotter) Enumeration(domain string, subdomains chan<- string) {
	urlAddress := fmt.Sprintf(c.Url+"/certs?domain=%s", domain)
	resp, err := http.Get(urlAddress)
	if err != nil {
		return
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
			subdomains <- subdomain.Hostname()
		}
	}

}
