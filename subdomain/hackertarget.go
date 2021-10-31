package subdomain

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ghaini/treasure-finder/constants"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HackerTarget struct {
	Url string
}

func NewHackerTarget() SubdomainFinderInterface {
	return &HackerTarget{
		Url: constants.HackertargetUrl,
	}
}

func (h HackerTarget) IsPaidProvider() bool {
	return false
}

func (h HackerTarget) Name() string {
	return "hackertarget"
}

func (h HackerTarget) SetAuth(token string) {
	return
}

func (h HackerTarget) GetAuth() string {return ""}

func (h HackerTarget) Enumeration(domain string) (result map[string]struct{}, statusCode int, err error) {
	result = make(map[string]struct{})
	urlAddress := fmt.Sprintf(h.Url+"/hostsearch/?q=%s", domain)
	client := http.Client{
		Timeout: 5 * time.Minute,
	}
	res, err := client.Get(urlAddress)
	if err != nil {
		return result, 500, err
	}
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, res.StatusCode, err
	}
	res.Body.Close()

	sc := bufio.NewScanner(bytes.NewReader(raw))
	for sc.Scan() {
		address := strings.SplitN(sc.Text(), ",", 2)
		if len(address) != 2 {
			continue
		}
		subdomain, err := url.Parse("http://" + address[0])
		if err != nil {
			continue
		}
		result[subdomain.Hostname()] = struct{}{}
	}

	return result, res.StatusCode, nil

}
