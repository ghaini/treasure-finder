package program

import (
	"regexp"
	"strings"
)

type ProgramFinderInterface interface {
	ProgramsList() ([]Program, error)
}

type ProgramFinder struct {
	Finders  []ProgramFinderInterface
	withIP   bool
	onlyStar bool
}

type Program struct {
	Name             string  `json:"name"`
	Provider         string  `json:"provider"`
	InScopeAssets    []Asset `json:"in_scope_assets"`
	OutOfScopeAssets []Asset `json:"out_of_scope_assets"`
}

type Asset struct {
	Address string `json:"address"`
	Bounty  bool   `json:"bounty"`
}

func NewProgramFinder() *ProgramFinder {
	return &ProgramFinder{
		Finders: []ProgramFinderInterface{
			NewHackerOne(),
			NewBugCrowd(),
		},
	}
}

func (p *ProgramFinder) GetPrograms() ([]Program, error) {
	var programs []Program
	programsMap := make( map[string]Program)
	checkIsIP, err := regexp.Compile("^\\d+\\.\\d+")
	if err != nil {
		return nil, err
	}

	for _, finder := range p.Finders {
		program, err := finder.ProgramsList()
		if err != nil {
			return nil, err
		}

		for _, pr := range program {
			var newInScopeAssets []Asset
			var newOutOfScopeAssets []Asset
			for _, asset := range pr.InScopeAssets {
				asset.Address = strings.TrimSpace(asset.Address)
				asset.Address = strings.ToLower(asset.Address)
				asset.Address = strings.TrimLeft(asset.Address, "https://")
				asset.Address = strings.TrimLeft(asset.Address, "http://")

				if !p.withIP && checkIsIP.MatchString(asset.Address) {
					continue
				}

				if p.onlyStar && !strings.HasPrefix(asset.Address, "*") {
					continue
				}

				if strings.Count(asset.Address, ".") > 4 {
					continue
				}

				newInScopeAssets = append(newInScopeAssets, asset)
			}

			for _, asset := range pr.OutOfScopeAssets {
				asset.Address = strings.TrimSpace(asset.Address)
				asset.Address = strings.ToLower(asset.Address)
				asset.Address = strings.TrimLeft(asset.Address, "https://")
				asset.Address = strings.TrimLeft(asset.Address, "http://")

				if !p.withIP && checkIsIP.MatchString(asset.Address) {
					continue
				}

				if strings.Count(asset.Address, ".") > 4 {
					continue
				}

				newOutOfScopeAssets = append(newOutOfScopeAssets, asset)
			}

			pr.InScopeAssets = newInScopeAssets
			pr.OutOfScopeAssets = newOutOfScopeAssets
			pr.Name = strings.ToLower(pr.Name)
			pr.Name = strings.ReplaceAll(pr.Name, " ", "-")
			if len(newInScopeAssets) > 0 {
				programsMap[pr.Name] = pr
			}
		}
	}

	for _, v := range programsMap{
		programs = append(programs, v)
	}

	return programs, nil
}

func (p *ProgramFinder) WithIP() {
	p.withIP = true
}

func (p *ProgramFinder) OnlyStar() {
	p.onlyStar = true
}
