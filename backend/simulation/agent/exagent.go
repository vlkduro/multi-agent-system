package agent

import (
	"math/rand/v2"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/simulation/environment"
)

// ExAgent est un exemple d'agent
type ExAgent struct {
	agent
	value    int
	decision int
}

type ExAgentJson struct {
	ID       string `json:"id"`
	Value    int    `json:"value"`
	Decision int    `json:"decision"`
}

func NewExAgent(id string, env *envpkg.Environment, syncChan chan bool) *ExAgent {
	exagt := &ExAgent{}
	exagt.agent = agent{
		iagt:     exagt,
		id:       envpkg.AgentID(id),
		env:      env,
		syncChan: syncChan,
	}

	return exagt
}

func (*ExAgent) Percept() {
	//TODO
}

func (agt *ExAgent) Deliberate() {
	factor := 1
	if rand.IntN(2) == 0 {
		factor = -1
	}
	agt.decision = factor * rand.IntN(100)
}

func (agt *ExAgent) Act() {
	agt.value += agt.decision
}

func (agt ExAgent) ToJsonObj() interface{} {
	return ExAgentJson{ID: string(agt.id), Value: agt.value, Decision: agt.decision}
}
