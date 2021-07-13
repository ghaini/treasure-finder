package treasureFinder

import (
	"github.com/ghaini/treasure-finder/program"
	"github.com/ghaini/treasure-finder/subdomain"
	"math/rand"
	"time"
)

type TreasureFinder struct {
	ProgramFinder   *program.ProgramFinder
	SubdomainFinder *subdomain.SubdomainFinder
}

func NewTreasureFinder() *TreasureFinder {
	rand.Seed(time.Now().UnixNano())
	return &TreasureFinder{
		ProgramFinder:   program.NewProgramFinder(),
		SubdomainFinder: subdomain.NewSubdomainFinder(),
	}
}
