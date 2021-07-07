package subdomain

import (
	"strings"
	"sync"
)

type SubdomainFinderInterface interface {
	Enumeration(domain string, subdomains chan<- string)
}

type SubdomainFinder struct {
	Finders []SubdomainFinderInterface
}

func NewSubdomainFinder() *SubdomainFinder {
	return &SubdomainFinder{
		[]SubdomainFinderInterface{
			NewUrlScan(),
			NewHackerTarget(),
			NewCrt(),
			NewThreatcrowd(),
			NewBufferover(),
			NewCertspotter(),
		},
	}
}

func (r SubdomainFinder) Enumeration(domain string) ([]string, error) {
	subdomainsUnionMap := make(map[string]struct{})
	wg := &sync.WaitGroup{}
	wg.Add(len(r.Finders))
	subdomainsChan := make(chan string)
	for _, finder := range r.Finders {
		go func(wg *sync.WaitGroup) {
			finder.Enumeration(domain, subdomainsChan)
			wg.Done()
		}(wg)
	}

	for subdomain := range subdomainsChan {
		subdomainsUnionMap[subdomain] = struct{}{}
	}

	wg.Wait()
	var subdomains []string
	for k, _ := range subdomainsUnionMap {
		if !strings.Contains(k, domain) {
			continue
		}

		subdomains = append(subdomains, k)
	}

	return subdomains, nil
}
