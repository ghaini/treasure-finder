package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
	"time"
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

func (c Certspotter) GetAuth() string {return ""}

func (c Certspotter) Name() string {
	return "certpotter"
}

func (c Certspotter) Enumeration(domain string) (result map[string]struct{}, statusCode int, err error) {
	result = make(map[string]struct{})
	urlAddress := fmt.Sprintf(c.Url+"/certs?domain=%s", domain)
	client := http.Client{
		Timeout: 5 * time.Minute,
	}
	resp, err := client.Get(urlAddress)
	if err != nil {
		return result, 500, err
	}
	defer resp.Body.Close()

	var certspotterResp []certspotterResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&certspotterResp)
	if err != nil {
		return result, resp.StatusCode, err
	}

	for _, csr := range certspotterResp {
		for _, address := range csr.DNSNames {
			subdomain, err := url.Parse("https://" + address)
			if err != nil {
				continue
			}
			result[subdomain.Hostname()] = struct{}{}
		}
	}

	return result, resp.StatusCode,nil
}
