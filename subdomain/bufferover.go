package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
	"strings"
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

func (b Bufferover) Enumeration(domain string, subdomains chan<- string) {
	urlAddress := fmt.Sprintf(b.Url+"?q=.%s", domain)
	resp, err := http.Get(urlAddress)
	if err != nil {
		return
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
		subdomain, err := url.Parse(subdomainAddress[0])
		if err != nil {
			continue
		}
		subdomains <- subdomain.Hostname()
	}
}
