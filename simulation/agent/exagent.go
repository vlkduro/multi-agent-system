package agent

import (
	"fmt"

	envpkg "../environment"
)

// ExAgent est un exemple d'agent
type ExAgent struct {
	id       envpkg.AgentID
	env      *envpkg.Environment
	syncChan chan bool
}

func NewExAgent(id string, env *envpkg.Environment, syncChan chan bool) *ExAgent {
	return &ExAgent{envpkg.AgentID(id), env, syncChan}
}

// L'agent est lancé en tant que microservice
func (agt *ExAgent) Start() {
	fmt.Printf("[%s] Starting agent\n", agt.ID())
	go func() {
		for {
			<-agt.syncChan
			agt.Percept()
			agt.Deliberate()
			agt.Act()
			agt.syncChan <- true
		}
	}()
}
func (*ExAgent) Percept() {
	//TODO
}
func (*ExAgent) Deliberate() {
	//TODO
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
