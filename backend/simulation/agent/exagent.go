package agent

import (
	"fmt"
	"math/rand/v2"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/simulation/environment"
)

// ExAgent est un exemple d'agent
type ExAgent struct {
	id       envpkg.AgentID
	env      *envpkg.Environment
	value    int
	syncChan chan bool
}

type ExAgentJson struct {
	ID    string `json:"id"`
	Value int    `json:"value"`
}

func NewExAgent(id string, env *envpkg.Environment, syncChan chan bool) *ExAgent {
	return &ExAgent{id: envpkg.AgentID(id), env: env, syncChan: syncChan}
}

// L'agent est lancé en tant que microservice
func (agt *ExAgent) Start() {
	fmt.Printf("[%s] Starting agent\n", agt.ID())
	go func() {
		for {
			run := <-agt.syncChan
			if !run {
				break
			}
			agt.Percept()
			agt.Deliberate()
			agt.Act()
			agt.syncChan <- run
		}
		fmt.Printf("[%s] Stopping agent\n", agt.ID())
	}()
}

func (*ExAgent) Percept() {
	//TODO
}

func (agt *ExAgent) Deliberate() {
	factor := 1
	if rand.IntN(2) == 0 {
		factor = -1
	}
	agt.value += factor * rand.IntN(100)
}

func (*ExAgent) Act() {
	//TODO
}

func (agt ExAgent) ID() envpkg.AgentID {
	return agt.id
}

func (agt ExAgent) GetSyncChan() chan bool {
	return agt.syncChan
}

func (agt ExAgent) ToJsonObj() interface{} {
	return ExAgentJson{ID: string(agt.id), Value: agt.value}
}
