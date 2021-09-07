package subdomain

import (
	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
)

type SubdomainFinderInterface interface {
	IsPaidProvider() bool
	SetAuth(token string)
	GetAuth() string
	Name() string
	Enumeration(domain string) (result map[string]struct{}, statusCode int, err error)
}

type SubdomainFinder struct {
	tokensPath string
	Finders    []SubdomainFinderInterface
}

type providerAuth struct {
	Tokens []providerAuthDetail `mapstructure:"Tokens"`
}

type providerAuthDetail struct {
	Token       string `mapstructure:"Token"`
	IsAvailable bool   `mapstructure:"IsAvailable"`
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
			NewCensys(),
			NewOmnisint(),
			NewRapiddns(),
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
		if finder.IsPaidProvider() && (r.tokensPath == "" || finder.GetAuth() == "")  {
			continue
		}

		wg.Add(1)
		go func(finder SubdomainFinderInterface, wg *sync.WaitGroup) {
			defer wg.Done()
			subdomainsMap, statusCode, err := finder.Enumeration(domain)
			if finder.IsPaidProvider() && (statusCode == 401 || statusCode == 403) {
				changeErr := r.changeTokenToUnavailable(finder.GetAuth(), finder.Name())
				if changeErr != nil {
					log.Println(changeErr)
				}
			}

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
		if !strings.HasSuffix(k, "."+domain) ||
			strings.Contains(k, "www.google.com") ||
			strings.Contains(k, "cloudflare.net") ||
			strings.Contains(k, "webproxy") ||
			strings.HasPrefix(k, "bvr") ||
			len(k) > (50+len(domain)) ||
			strings.Count(k, "-") > 3 {
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

		var authProvider providerAuth
		err = viperInstance.Unmarshal(&authProvider)
		if err != nil {
			log.Println(err)
			continue
		}

		for i := range authProvider.Tokens {
			j := rand.Intn(i + 1)
			authProvider.Tokens[i], authProvider.Tokens[j] = authProvider.Tokens[j], authProvider.Tokens[i]
		}

		for _, token := range authProvider.Tokens {
			if token.IsAvailable {
				finder.SetAuth(token.Token)
				break
			}
		}
	}
}

func (r *SubdomainFinder) changeTokenToUnavailable(token, provider string) error {
	tokenPath := r.tokensPath + "/" + provider + ".toml"
	viperInstance := viper.New()
	viperInstance.SetConfigFile(tokenPath)
	err := viperInstance.ReadInConfig()
	if err != nil {
		return err
	}

	var authProvider providerAuth
	err = viperInstance.Unmarshal(&authProvider)
	if err != nil {
		return err
	}

	for i, auth := range authProvider.Tokens {
		if auth.Token == token {
			f, err := os.Create(tokenPath)
			if err != nil {
				return err
			}

			authProvider.Tokens[i].IsAvailable = false
			if err = toml.NewEncoder(f).Encode(authProvider); err != nil {
				return err
			}

			if err = f.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
