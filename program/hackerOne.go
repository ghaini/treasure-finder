package program

import (
	"encoding/json"
	"github.com/ghaini/treasure-finder/constants"
	"io/ioutil"
	"net/http"
)

type HackerOne struct {
	Url string
}

type HackerOneResponse struct {
	Name    string `json:"name"`
	Targets struct {
		InScope []struct {
			AssetIdentifier   string `json:"asset_identifier"`
			EligibleForBounty bool   `json:"eligible_for_bounty"`
		} `json:"in_scope"`
		OutOFScope []struct {
			AssetIdentifier   string `json:"asset_identifier"`
			EligibleForBounty bool   `json:"eligible_for_bounty"`
		} `json:"out_of_scope,omitempty"`
	} `json:"targets"`
}

func NewHackerOne() ProgramFinderInterface {
	return HackerOne{Url: constants.HackerOneFileUrl}
}

func (h HackerOne) ProgramsList() ([]Program, error) {
	resp, err := http.Get(h.Url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var hackerOneResponses []HackerOneResponse
	err = json.Unmarshal(body, &hackerOneResponses)
	if err != nil {
		return nil, err
	}

	var programs []Program
	for _, hackerOneResponse := range hackerOneResponses {
		var inScopeAssets []Asset
		var outOfScopeAssets []Asset
		for _, target := range hackerOneResponse.Targets.InScope {
			asset := Asset{
				Address: target.AssetIdentifier,
				Bounty:  target.EligibleForBounty,
			}
			inScopeAssets = append(inScopeAssets, asset)
		}

		for _, target := range hackerOneResponse.Targets.OutOFScope {
			asset := Asset{
				Address: target.AssetIdentifier,
				Bounty:  target.EligibleForBounty,
			}
			outOfScopeAssets = append(outOfScopeAssets, asset)
		}

		programs = append(programs, Program{
			Name:             hackerOneResponse.Name,
			InScopeAssets:    inScopeAssets,
			OutOfScopeAssets: outOfScopeAssets,
		})
	}

	return programs, nil
}
