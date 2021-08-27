package subdomain

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/ghaini/treasure-finder/constants"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Censys struct {
	Url    string
	ApiKey string
}

type censysResponse struct {
	Results  []censysResult `json:"results"`
	Metadata censysMetadata `json:"metadata"`
}

type censysResult struct {
	Names []string `json:"parsed.names"`
}

type censysMetadata struct {
	Count int `json:"count"`
	Page  int `json:"page"`
	Pages int `json:"pages"`
}

type censysRequest struct {
	Query  string   `json:"query"`
	Fields []string `json:"fields"`
	Page   int      `json:"page"`
}

func NewCensys() SubdomainFinderInterface {
	return &Censys{
		Url: constants.CensysUrl,
	}
}

func (c Censys) IsPaidProvider() bool {
	return true
}

func (c Censys) Name() string {
	return "censys"
}

func (c Censys) SetAuth(token string) {
	c.ApiKey = token
	return
}

func (c Censys) GetAuth() string {
	return c.ApiKey
}

func (c Censys) Enumeration(domain string) (map[string]struct{}, int, error) {
	resultMap := make(map[string]struct{})
	page := 1
	maxPage := 1
	for {
		censysResponse, statusCode, err := c.censysRequest(domain, page)
		if err != nil {
			return nil, statusCode, err
		}
		maxPage = censysResponse.Metadata.Pages
		for _, result := range censysResponse.Results {
			for _, name := range result.Names {
				subdomain, err := url.Parse("https://" + name)
				if err != nil {
					continue
				}
				resultMap[subdomain.Hostname()] = struct{}{}
			}
		}

		if maxPage == page {
			break
		}
		page++
	}

	return resultMap, 200, nil
}

func (c Censys) censysRequest(domain string, page int) (*censysResponse, int, error) {
	req := censysRequest{
		Query:  domain,
		Fields: []string{"parsed.names"},
		Page:   page,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return nil, 500, err
	}

	client := &http.Client{}
	httpReq, err := http.NewRequest("POST", constants.CensysUrl, bytes.NewBuffer(reqJson))
	if err != nil {
		return nil, 500, err
	}
	httpReq.Header.Set("Authorization", base64.StdEncoding.EncodeToString([]byte(c.ApiKey)))
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, 500, err
	}

	defer resp.Body.Close()

	var censysResponse *censysResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	err = json.Unmarshal(body, &censysResponse)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return censysResponse, resp.StatusCode, nil
}
