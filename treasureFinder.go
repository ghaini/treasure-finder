package treasureFinder

import (
	"github.com/ghaini/treasure-finder/program"
	"github.com/ghaini/treasure-finder/subdomain"
)

type TreasureFinder struct {
	programFinder   *program.ProgramFinder
	subdomainFinder *subdomain.SubdomainFinder
}

func NewTreasureFinder() *TreasureFinder {
	return &TreasureFinder{
		programFinder:   program.NewProgramFinder(),
		subdomainFinder: subdomain.NewSubdomainFinder(),
	}
}
