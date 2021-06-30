package program

import (
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
	for _, finder := range p.Finders {
		program, err := finder.ProgramsList()
		if err != nil {
			return nil, err
		}

		for _, pr := range program{
			var newAssetsIdentifier []string
			for _, asset := range pr.AssetsIdentifier  {
				if !p.withIP && strings.Count(".", asset) == 3 {
					continue
				}

				if p.onlyStar && !strings.Contains(asset, "*") {
					continue
				}

				newAssetsIdentifier = append(newAssetsIdentifier, asset)
			}
			pr.AssetsIdentifier = newAssetsIdentifier

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
