package agent

import (
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/simulation/environment"
	"time"
)

type job int

const (
	Worker = iota
	Guardian
	Forager
)

// BeeAgent hérite de /simulation/agent/agent.go "struct Agent"
// Interface IAgent
type BeeAgent struct {
	agent
	ruche     Ruche
	birthDate time.Time
	maxNectar int
	job       int
}

type BeeAgentJson struct {
	ID        string `json:"id"`
	MaxNectar int    `json:"maxNectar"`
	Job       int    `json:"job"`
}

func NewBeeAgent(id string, env *envpkg.Environment, syncChan chan bool, s int, r Ruche, bd time.Time, maxnectar int, job int) *BeeAgent {
	beeAgent := &BeeAgent{}
	beeAgent.agent = agent{
		iagt:     beeAgent,
		id:       envpkg.AgentID(id),
		env:      env,
		syncChan: syncChan,
		speed:    s,
	}
	beeAgent.ruche = r
	beeAgent.birthDate = bd
	beeAgent.maxNectar = maxnectar
	beeAgent.job = job
	return beeAgent
}

func (agt *BeeAgent) Start() {
}

func (agt *BeeAgent) Stop() {

}

func (agt *BeeAgent) Percept() {
}

func (agt *BeeAgent) Deliberate() {
}

func (agt *BeeAgent) Act() {
}

func (agt *BeeAgent) ToJsonObj() interface{} {
	return BeeAgentJson{ID: string(agt.id), MaxNectar: agt.maxNectar, Job: agt.job}
}
