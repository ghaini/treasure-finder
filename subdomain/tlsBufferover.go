package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"net/http"
	"net/url"
	"strings"
)

type TLSBufferover struct {
	Url string
}

type TLSbufferoverResponse struct {
	Records []string `json:"FDNS_A"`
}

func NewTLSBufferover() SubdomainFinderInterface {
	return &TLSBufferover{
		Url: constants.TLSBufferoverUrl,
	}
}

func (b TLSBufferover) IsPaidProvider() bool {
	return false
}

func (b TLSBufferover) Name() string {
	return "bufferover"
}

func (b TLSBufferover) SetAuth(token string) {
	return
}

func (b TLSBufferover) GetAuth() string {return ""}

func (b TLSBufferover) Enumeration(domain string) (result map[string]struct{}, statusCode int, err error) {
	urlAddress := fmt.Sprintf(b.Url+"?q=%s", domain)
	resp, err := http.Get(urlAddress)
	if err != nil {
		return result, 500,err
	}
	defer resp.Body.Close()

	var bufferoverRes bufferoverResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&bufferoverRes)
	if err != nil {
		return result, resp.StatusCode, err
	}

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
