package subdomain

import "strings"

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
		},
	}
}

func (r SubdomainFinder) Enumeration(domain string) ([]string, error) {
	subdomainsUnionMap := make(map[string]struct{})
	for _, finder := range r.Finders {
		subdomainMap, err := finder.Enumeration(domain)
		if err != nil {
			return nil, err
		}

		for k, v := range subdomainMap {
			subdomainsUnionMap[k] = v
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
