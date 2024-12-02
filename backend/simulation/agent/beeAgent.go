package agent

import (
	env "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/simulation/environment"
	"time"
)

type job int

const (
	Worker = iota
	Guardian
	Forager
)

type BeeAgent struct {
	ruche     Ruche
	birthDate time.Time
	maxNectar int
	job       int
}

func NewBeeAgent(r Ruche, bd time.Time, maxnectar int, job int) *BeeAgent {
	return &BeeAgent{r, bd, maxnectar, job}
}

func (agt *BeeAgent) Start() env.IAgent {

}

func (agt *BeeAgent) Stop() {

}
