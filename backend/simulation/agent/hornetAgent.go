package agent

import (
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/simulation/environment"
	"time"
)

// HornetAgent hérite de /simulation/agent/agent.go "struct Agent"
// Interface IAgent
type HornetAgent struct {
	agent
	ruche     Ruche
	birthDate time.Time
	maxNectar int
	job       int
}

type HornetAgentJson struct {
	ID string `json:"id"`
}

func NewHornetAgent(id string, env *envpkg.Environment, syncChan chan bool, s int) *HornetAgent {
	hAgent := &HornetAgent{}
	hAgent.agent = agent{
		iagt:     hAgent,
		id:       envpkg.AgentID(id),
		env:      env,
		syncChan: syncChan,
		speed:    s,
	}

	return hAgent
}

func (agt *HornetAgent) Start() {
}

func (agt *HornetAgent) Stop() {

}

func (agt *HornetAgent) Percept() {
}

func (agt *HornetAgent) Deliberate() {
}

func (agt *HornetAgent) Act() {
}

func (agt *HornetAgent) ToJsonObj() interface{} {
	return BeeAgentJson{ID: string(agt.id)}
}
