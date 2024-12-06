package agent

import (
	"fmt"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

// Abstract class used with a template pattern
// - iagt: An interface representing the agent's actions
// - id: An identifier for the agent
// - env: A pointer to the environment in which the agent operates
// - syncChan: A channel used for synchronization purposes
type Agent struct {
	iagt     envpkg.IAgent
	id       envpkg.AgentID
	pos      *envpkg.Position
	env      *envpkg.Environment
	syncChan chan bool
	Speed    int
}

// Agent is launched as a microservice
func (agt *Agent) Start() {
	fmt.Printf("[%s] Starting Agent\n", agt.ID())
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
		fmt.Printf("[%s] Stopping Agent\n", agt.ID())
	}()
}

func (agt Agent) ID() envpkg.AgentID {
	return agt.id
}

func (agt Agent) GetSyncChan() chan bool {
	return agt.syncChan
}

func (agt Agent) Position() *envpkg.Position {
	if agt.pos == nil {
		return nil
	}
	return agt.pos.Copy()
}
