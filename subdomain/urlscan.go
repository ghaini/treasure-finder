package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
)

type UrlScan struct {
	Url string
}

type urlScanResponse struct {
	Results []struct {
		Task struct {
			URL string `json:"url"`
		} `json:"task"`

		Page struct {
			URL string `json:"url"`
		} `json:"page"`
	} `json:"results"`
}

func NewUrlScan() SubdomainFinderInterface {
	return &UrlScan{
		Url: constants.UrlScanUrl,
	}
}

func (u UrlScan) Enumeration(domain string, subdomains chan<- string) {
	fetchURL := fmt.Sprintf(u.Url+"/search/?q=domain:%s", domain)
	resp, err := http.Get(fetchURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	var urlScanResponse urlScanResponse
	err = dec.Decode(&urlScanResponse)
	for _, r := range urlScanResponse.Results {
		subdomain, err := url.Parse(r.Task.URL)
		if err != nil {
			continue
		}
		subdomains <- subdomain.Hostname()
	}
}
