package agent

import (
	"time"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
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

func NewBeeAgent(id string, env *envpkg.Environment, syncChan chan bool, speed int, r Ruche, dob time.Time, maxnectar int, job int) *BeeAgent {
	beeAgent := &BeeAgent{}
	beeAgent.agent = agent{
		iagt:     beeAgent,
		id:       envpkg.AgentID(id),
		env:      env,
		syncChan: syncChan,
		speed:    speed,
	}
	beeAgent.ruche = r
	beeAgent.birthDate = dob
	beeAgent.maxNectar = maxnectar
	beeAgent.job = job
	return beeAgent
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
