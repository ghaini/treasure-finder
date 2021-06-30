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
			AssetIdentifier string `json:"asset_identifier"`
		} `json:"in_scope"`
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
		var assetsIdentifier []string
		for _, target := range hackerOneResponse.Targets.InScope {
			assetsIdentifier = append(assetsIdentifier, target.AssetIdentifier)
		}

		programs = append(programs, Program{
			Name:             hackerOneResponse.Name,
			AssetsIdentifier: assetsIdentifier,
		})
	}

	return programs, nil
}
