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
)

type HackerTarget struct {
	Url string
}

func NewHackerTarget() SubdomainFinderInterface {
	return &HackerTarget{
		Url: constants.HackertargetUrl,
	}
}

func (h HackerTarget) Enumeration(domain string, subdomains chan<- string) {
	urlAddress := fmt.Sprintf(h.Url+"/hostsearch/?q=%s", domain)
	res, err := http.Get(urlAddress)
	if err != nil {
		return
	}
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	res.Body.Close()

	sc := bufio.NewScanner(bytes.NewReader(raw))
	for sc.Scan() {
		address := strings.SplitN(sc.Text(), ",", 2)
		if len(address) != 2 {
			continue
		}
		subdomain, err := url.Parse(address[0])
		if err != nil {
			continue
		}
		subdomains <- subdomain.Hostname()
	}

}
