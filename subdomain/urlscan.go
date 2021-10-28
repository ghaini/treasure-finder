package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
	"time"
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

func (u UrlScan) IsPaidProvider() bool {
	return false
}

func (u UrlScan) Name() string {
	return "urlscan"
}

func (u UrlScan) SetAuth(token string) {
	return
}

func (u UrlScan) GetAuth() string { return "" }

func (u UrlScan) Enumeration(domain string) (result map[string]struct{}, statusCode int, err error) {
	result = make(map[string]struct{})
	fetchURL := fmt.Sprintf(u.Url+"/search/?q=domain:%s", domain)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(fetchURL)
	if err != nil {
		return result, 500, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)

	var urlScanResponse urlScanResponse
	err = dec.Decode(&urlScanResponse)
	if err != nil {
		return result, resp.StatusCode, err
	}

	for _, r := range urlScanResponse.Results {
		subdomain, err := url.Parse(r.Task.URL)
		if err != nil {
			continue
		}
		result[subdomain.Hostname()] = struct{}{}
	}

	return result, resp.StatusCode, nil
}
