package program

import (
	"encoding/json"
	"github.com/ghaini/treasure-finder/constants"
	"io/ioutil"
	"net/http"
)

type BugCrowd struct {
	Url string
}

type BugCrowdResponse struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Targets struct {
		InScope []struct {
			Target string `json:"target"`
		} `json:"in_scope"`
		OutOFScope []struct {
			Target string `json:"target"`
		} `json:"out_of_scope,omitempty"`
	} `json:"targets"`
}

func NewBugCrowd() ProgramFinderInterface {
	return BugCrowd{Url: constants.BugCrowdFileUrl}
}

func (b BugCrowd) ProgramsList() ([]Program, error) {
	resp, err := http.Get(b.Url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var bugCrowdResponses []BugCrowdResponse
	err = json.Unmarshal(body, &bugCrowdResponses)
	if err != nil {
		return nil, err
	}

	var programs []Program
	for _, bugCrowdResponse := range bugCrowdResponses {
		var inScopeAssets []Asset
		var outOfScopeAssets []Asset
		for _, target := range bugCrowdResponse.Targets.InScope {
			asset := Asset{
				Address: target.Target,
				Bounty:  true,
			}
			inScopeAssets = append(inScopeAssets, asset)
		}

		for _, target := range bugCrowdResponse.Targets.OutOFScope {
			asset := Asset{
				Address: target.Target,
				Bounty:  true,
			}
			outOfScopeAssets = append(outOfScopeAssets, asset)
		}

		programs = append(programs, Program{
			Name:             bugCrowdResponse.Name,
			URL:             bugCrowdResponse.URL,
			Provider:         "bugCrowd",
			InScopeAssets:    inScopeAssets,
			OutOfScopeAssets: outOfScopeAssets,
		})
	}

	return programs, nil
}
