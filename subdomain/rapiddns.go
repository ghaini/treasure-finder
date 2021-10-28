package subdomain

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ghaini/treasure-finder/constants"
)

type Rapiddns struct {
	Url string
}

func NewRapiddns() SubdomainFinderInterface {
	return &Rapiddns{
		Url: constants.RapiddnsUrl,
	}
}

func (r Rapiddns) IsPaidProvider() bool {
	return false
}

func (r Rapiddns) Name() string {
	return "rapiddns"
}

func (r Rapiddns) SetAuth(token string) {
	return
}

func (r Rapiddns) GetAuth() string {return ""}

func (r Rapiddns) Enumeration(domain string) (result map[string]struct{}, statusCode int, err error) {
	result = make(map[string]struct{})
	urlAddress := fmt.Sprintf(r.Url+"%s?full=1#result", domain)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(urlAddress)
	if err != nil {
		return result, 500, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return result, resp.StatusCode, err
	}

	doc.Find(".table td").Each(func(i int, s *goquery.Selection) {
		address := s.Text()
		if strings.Contains(address, domain) {
			subdomain, err := url.Parse("https://" + address)
			if err != nil {
				return
			}

			result[subdomain.Hostname()] = struct{}{}
		}
	})

	return result, resp.StatusCode, nil
}
