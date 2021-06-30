package treasureFinder

import (
	"github.com/ghaini/treasure-finder/program"
	"github.com/ghaini/treasure-finder/subdomain"
)

type TreasureFinder struct {
	ProgramFinder   *program.ProgramFinder
	SubdomainFinder *subdomain.SubdomainFinder
}

func NewTreasureFinder() *TreasureFinder {
	return &TreasureFinder{
		ProgramFinder:   program.NewProgramFinder(),
		SubdomainFinder: subdomain.NewSubdomainFinder(),
	}
}
