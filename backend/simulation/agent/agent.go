package agent

import (
	"fmt"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/simulation/environment"
)

// Abstract class used with a template pattern
// - iagt: An interface representing the agent's actions
// - id: An identifier for the agent
// - env: A pointer to the environment in which the agent operates
// - syncChan: A channel used for synchronization purposes
type agent struct {
	iagt     envpkg.IAgent
	id       envpkg.AgentID
	env      *envpkg.Environment
	syncChan chan bool
}

// Agent is launched as a microservice
func (agt *agent) Start() {
	fmt.Printf("[%s] Starting agent\n", agt.ID())
	go func() {
		for {
			run := <-agt.syncChan
			if !run {
				agt.syncChan <- run
				break
			}
			agt.iagt.Percept()
			agt.iagt.Deliberate()
			agt.iagt.Act()
			agt.syncChan <- run
		}
		fmt.Printf("[%s] Stopping agent\n", agt.ID())
	}()
}

func (agt agent) ID() envpkg.AgentID {
	return agt.id
}

func (agt agent) GetSyncChan() chan bool {
	return agt.syncChan
}
