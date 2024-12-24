package agent

import (
	"fmt"

	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/agent/vision"
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

// Abstract class used with a template pattern
// - iagt: An interface representing the agent's actions
// - id: An identifier for the agent
// - env: A pointer to the environment in which the agent operates
// - syncChan: A channel used for synchronization purposes
type Agent struct {
	iagt        envpkg.IAgent
	id          envpkg.AgentID
	pos         *envpkg.Position
	orientation envpkg.Orientation
	env         *envpkg.Environment
	visionFunc  vision.VisionFunc
	syncChan    chan bool
	Speed       int
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

func (agt Agent) see() []vision.SeenElem {
	return agt.visionFunc(agt.iagt, agt.env)
}

func (agt Agent) Orientation() envpkg.Orientation {
	return agt.orientation
}

func (agt *Agent) gotoNextStepTowards(pos *envpkg.Position) {
	// Not optimized,in some cases the agent will go in a direction that is not the best
	// But it is a simple way to make the agent move
	//  | a | o | a os agent and o is object
	//  |   | X | X is destination
	//  Best scenario   Worst scenario
	//  | a | o | SE   | a>| o | ES
	//  | v | > |      | v |   |
	for i := 0; i < agt.Speed; i++ {
		if agt.pos.X < pos.X {
			agt.pos.GoDown(agt.env.GetMap())
			if agt.pos.Y < pos.Y {
				agt.pos.GoRight(agt.env.GetMap())
				agt.orientation = envpkg.SouthEast
			} else if agt.pos.Y > pos.Y {
				agt.pos.GoLeft(agt.env.GetMap())
				agt.orientation = envpkg.SouthWest
			} else {
				agt.orientation = envpkg.South
			}
		} else if agt.pos.X > pos.X {
			agt.pos.GoUp(agt.env.GetMap())
			if agt.pos.Y < pos.Y {
				agt.pos.GoRight(agt.env.GetMap())
				agt.orientation = envpkg.NorthEast
			} else if agt.pos.Y > pos.Y {
				agt.pos.GoLeft(agt.env.GetMap())
				agt.orientation = envpkg.NorthWest
			} else {
				agt.orientation = envpkg.North
			}
		} else {
			if agt.pos.Y < pos.Y {
				agt.pos.GoRight(agt.env.GetMap())
				agt.orientation = envpkg.East
			} else if agt.pos.Y > pos.Y {
				agt.pos.GoLeft(agt.env.GetMap())
				agt.orientation = envpkg.West
			}
		}
	}

}
