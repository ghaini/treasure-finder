package subdomain

import (
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"strings"
	"sync"
)

type SubdomainFinderInterface interface {
	IsPaidProvider() bool
	SetAuth(token string)
	Name() string
	Enumeration(domain string) (map[string]struct{}, error)
}

type SubdomainFinder struct {
	tokensPath string
	Finders    []SubdomainFinderInterface
}

type providerAuth struct {
	Tokens []string `mapstructure:"tokens"`
}

func NewSubdomainFinder() *SubdomainFinder {
	return &SubdomainFinder{
		Finders: []SubdomainFinderInterface{
			NewUrlScan(),
			NewHackerTarget(),
			NewCrt(),
			NewThreatcrowd(),
			NewCertspotter(),
			NewBufferover(),
			NewSecuritytrails(),
			NewTLSBufferover(),
			NewJLDC(),
			NewAlienvault(),
		},
	}
}

func (r *SubdomainFinder) Enumeration(domain string) ([]string, error) {

	if r.tokensPath != "" {
		r.initialPaidProviders(r.tokensPath)
	}

	subdomainsUnionMap := make(map[string]struct{})
	wg := &sync.WaitGroup{}
	subdomainsMapChan := make(chan map[string]struct{})
	for _, finder := range r.Finders {
		if r.tokensPath == "" && finder.IsPaidProvider() {
			continue
		}

		wg.Add(1)
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
		for subdomain, _ := range subdomainsMap {
			subdomainsUnionMap[subdomain] = struct{}{}
		}
	}

	var subdomains []string
	for k, _ := range subdomainsUnionMap {
		k = strings.ReplaceAll(k, "https://", "")
		k = strings.ReplaceAll(k, "http://", "")
		k = strings.ReplaceAll(k, "*", "")
		k = strings.Trim(k, ".")
		k = strings.Trim(k, "/")
		k = strings.TrimSpace(k)
		if !strings.Contains(k, "." + domain) ||
			strings.Contains(k, "www.google.com") ||
			strings.Contains(k, "webproxy") ||
			strings.HasPrefix(k, "bvr") ||
			len(k) > 50 {
			continue
		}

		subdomains = append(subdomains, k)
	}

	subdomains = append(subdomains, domain)
	return subdomains, nil
}

func (r *SubdomainFinder) SetUsePaidProviders(baseTokensPath string) {
	r.tokensPath = baseTokensPath
}

func (r *SubdomainFinder) initialPaidProviders(baseTokenPath string) {
	for _, finder := range r.Finders {
		if !finder.IsPaidProvider() {
			continue
		}

		tokenPath := baseTokenPath + "/" + finder.Name() + ".toml"
		viperInstance := viper.New()
		viperInstance.SetConfigFile(tokenPath)
		err := viperInstance.ReadInConfig()
		if err != nil {
			log.Println(err)
			continue
		}

		var auth providerAuth
		err = viperInstance.Unmarshal(&auth)
		if err != nil {
			log.Println(err)
			continue
		}

		token := auth.Tokens[rand.Intn(len(auth.Tokens))]
		finder.SetAuth(token)
	}
}
