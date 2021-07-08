package subdomain

import (
	"strings"
	"sync"
)

type SubdomainFinderInterface interface {
	Enumeration(domain string) (map[string]struct{}, error)
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
			NewCertspotter(),
			NewBufferover(),
		},
	}
}


func (r SubdomainFinder) Enumeration(domain string) ([]string, error) {
	subdomainsUnionMap := make(map[string]struct{})
	wg := &sync.WaitGroup{}
	wg.Add(len(r.Finders))
	subdomainsMapChan := make(chan map[string]struct{})
	for _, finder := range r.Finders {
		go func(finder SubdomainFinderInterface, wg *sync.WaitGroup) {
			defer wg.Done()
			subdomainsMap, err := finder.Enumeration(domain)
			if err != nil {
				return
			}
			subdomainsMapChan <- subdomainsMap
		}(finder, wg)
	}

	go func() {
		wg.Wait()
		close(subdomainsMapChan)
	}()

	for subdomainsMap := range subdomainsMapChan {
		for subdomain, _ := range subdomainsMap{
			subdomainsUnionMap[subdomain] = struct{}{}
		}
	}

	var subdomains []string
	for k, _ := range subdomainsUnionMap {
		if !strings.Contains(k, domain) {
			continue
		}

		subdomains = append(subdomains, k)
	}

	return subdomains, nil
}