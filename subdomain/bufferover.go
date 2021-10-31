package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Bufferover struct {
	Url string
}

type bufferoverResponse struct {
	Records []string `json:"FDNS_A"`
}

func NewBufferover() SubdomainFinderInterface {
	return &Bufferover{
		Url: constants.BufferoverUrl,
	}
}

func (b Bufferover) IsPaidProvider() bool {
	return false
}

func (b Bufferover) Name() string {
	return "bufferover"
}

func (b Bufferover) SetAuth(token string) {
	return
}

func (b Bufferover) GetAuth() string {return ""}

func (b Bufferover) Enumeration(domain string) (result map[string]struct{}, statusCode int, err error) {
	result = make(map[string]struct{})
	urlAddress := fmt.Sprintf(b.Url+"?q=.%s", domain)
	client := http.Client{
		Timeout: 5 * time.Minute,
	}
	resp, err := client.Get(urlAddress)
	if err != nil {
		return result, 500, err
	}
	defer resp.Body.Close()

	var bufferoverRes bufferoverResponse
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&bufferoverRes)
	for _, record := range bufferoverRes.Records {
		subdomainAddress := strings.SplitN(record, ",", 2)
		if len(subdomainAddress) != 2 {
			continue
		}
		subdomain, err := url.Parse("https://" + subdomainAddress[0])
		if err != nil {
			continue
		}
		result[subdomain.Hostname()] = struct{}{}
	}

	return result, resp.StatusCode, nil
}
