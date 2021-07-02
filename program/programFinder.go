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
	Name             string   `json:"name"`
	AssetsIdentifier []string `json:"assets_identifier"`
}

func NewProgramFinder() *ProgramFinder {
	return &ProgramFinder{
		Finders: []ProgramFinderInterface{
			NewHackerOne(),
		},
	}
}

func (p *ProgramFinder) GetPrograms() ([]Program, error) {
	var programs []Program
	checkIsIP, err := regexp.Compile("^\\d+\\.\\d+")
	if err != nil {
		return nil, err
	}
	
	for _, finder := range p.Finders {
		program, err := finder.ProgramsList()
		if err != nil {
			return nil, err
		}
		
		for _, pr := range program{
			var newAssetsIdentifier []string
			for _, asset := range pr.AssetsIdentifier  {
				asset = strings.TrimSpace(asset)
				asset = strings.ToLower(asset)
				asset = strings.TrimLeft(asset, "https://")
				asset = strings.TrimLeft(asset, "http://")

				if !p.withIP && checkIsIP.MatchString(asset) {
					continue
				}

				if p.onlyStar && !strings.HasPrefix(asset, "*") {
					continue
				}

				if strings.Count(asset, ".") > 4 {
					continue
				}
				
				newAssetsIdentifier = append(newAssetsIdentifier, asset)
			}

			pr.AssetsIdentifier = newAssetsIdentifier
			pr.Name = strings.ToLower(pr.Name)
			pr.Name = strings.ReplaceAll(pr.Name, " ", "-")
			if len(newAssetsIdentifier) > 0 {
				programs = append(programs, pr)
			}
		}
	}
	return programs, nil
}

func (p *ProgramFinder) WithIP() {
	p.withIP = true
}

func (p *ProgramFinder) OnlyStar() {
	p.onlyStar = true
}
