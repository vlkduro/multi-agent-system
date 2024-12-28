package agent

import (
	"fmt"
	"math/rand"

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

func (agt Agent) see() []*vision.SeenElem {
	return agt.visionFunc(agt.iagt, agt.env)
}

func (agt Agent) Orientation() envpkg.Orientation {
	return agt.orientation
}

func (agt *Agent) goNorth() bool {
	success := agt.pos.GoNorth(agt.env.GetMap())
	agt.orientation = envpkg.North
	return success
}

func (agt *Agent) goSouth() bool {
	success := agt.pos.GoSouth(agt.env.GetMap())
	agt.orientation = envpkg.South
	return success
}

func (agt *Agent) goEast() bool {
	success := agt.pos.GoEast(agt.env.GetMap())
	agt.orientation = envpkg.East
	return success
}

func (agt *Agent) goWest() bool {
	success := agt.pos.GoWest(agt.env.GetMap())
	agt.orientation = envpkg.West
	return success
}

func (agt *Agent) goNorthEast() bool {
	success := agt.pos.GoNorthEast(agt.env.GetMap())
	agt.orientation = envpkg.NorthEast
	return success
}

func (agt *Agent) goSouthEast() bool {
	success := agt.pos.GoSouthEast(agt.env.GetMap())
	agt.orientation = envpkg.SouthEast
	return success
}

func (agt *Agent) goNorthWest() bool {
	success := agt.pos.GoNorthWest(agt.env.GetMap())
	agt.orientation = envpkg.NorthWest
	return success
}

func (agt *Agent) goSouthWest() bool {
	success := agt.pos.GoSouthWest(agt.env.GetMap())
	agt.orientation = envpkg.SouthWest
	return success
}

// https://web.archive.org/web/20171022224528/http://www.policyalmanac.org:80/games/aStarTutorial.htm
func (agt *Agent) gotoNextStepTowards(pos *envpkg.Position) {
	chain := agt.env.PathFinding(agt.pos, pos)
	// We remove the first element who is the current position of the agent
	if len(chain) > 0 && chain[0].Equal(agt.pos) {
		chain = chain[1:]
	}
	fmt.Printf("\n[%s] Going to [%d %d] : [%d %d] -> ", agt.id, pos.X, pos.Y, agt.pos.X, agt.pos.Y)
	for i := 0; i < agt.Speed && i < len(chain); i++ {
		pos := chain[i]
		movementSuccess := false
		if agt.pos.X < pos.X {
			if agt.pos.Y < pos.Y {
				movementSuccess = agt.goSouthEast()
			} else if agt.pos.Y > pos.Y {
				movementSuccess = agt.goNorthEast()
			} else {
				movementSuccess = agt.goEast()
			}
		} else if agt.pos.X > pos.X {
			if agt.pos.Y < pos.Y {
				movementSuccess = agt.goSouthWest()
			} else if agt.pos.Y > pos.Y {
				movementSuccess = agt.goNorthWest()
			} else {
				movementSuccess = agt.goWest()
			}
		} else {
			if agt.pos.Y < pos.Y {
				movementSuccess = agt.goSouth()
			} else if agt.pos.Y > pos.Y {
				movementSuccess = agt.goNorth()
			}
		}
		// If the movement was not successful, we try to go in another direction
		if !movementSuccess {
		}
		// 	directions := []envpkg.Orientation{envpkg.North, envpkg.East, envpkg.South, envpkg.West}
		// 	rand.Shuffle(len(directions), func(i, j int) {
		// 		directions[i], directions[j] = directions[j], directions[i]
		// 	})
		// 	for _, orient := range directions {
		// 		if movementSuccess {
		// 			break
		// 		}
		// 		switch orient {
		// 		case envpkg.North:
		// 			movementSuccess = agt.goNorth()
		// 		case envpkg.South:
		// 			movementSuccess = agt.goSouth()
		// 		case envpkg.East:
		// 			movementSuccess = agt.goEast()
		// 		case envpkg.West:
		// 			movementSuccess = agt.goWest()
		// 		}
		// 	}
		// }
	}
}

func (agt *Agent) getNextWanderingPosition() *envpkg.Position {
	nextWanderingOrientation := agt.orientation
	// Chances : 3/4 th keeping the same orientation, 1/8th changing to the left, 1/8th changing to the right
	chancesToChangeOrientation := rand.Intn(8)
	// To the left
	if chancesToChangeOrientation < 2 {
		switch agt.orientation {
		case envpkg.North:
			if chancesToChangeOrientation == 0 {
				nextWanderingOrientation = envpkg.NorthWest
			} else {
				nextWanderingOrientation = envpkg.NorthEast
			}
		case envpkg.South:
			if chancesToChangeOrientation == 0 {
				nextWanderingOrientation = envpkg.SouthWest
			} else {
				nextWanderingOrientation = envpkg.SouthEast
			}
		case envpkg.East:
			if chancesToChangeOrientation == 0 {
				nextWanderingOrientation = envpkg.NorthEast
			} else {
				nextWanderingOrientation = envpkg.SouthEast
			}
		case envpkg.West:
			if chancesToChangeOrientation == 0 {
				nextWanderingOrientation = envpkg.NorthWest
			} else {
				nextWanderingOrientation = envpkg.SouthWest
			}
		case envpkg.NorthEast:
			if chancesToChangeOrientation == 0 {
				nextWanderingOrientation = envpkg.North
			} else {
				nextWanderingOrientation = envpkg.East
			}
		case envpkg.NorthWest:
			if chancesToChangeOrientation == 0 {
				nextWanderingOrientation = envpkg.North
			} else {
				nextWanderingOrientation = envpkg.West
			}
		case envpkg.SouthEast:
			if chancesToChangeOrientation == 0 {
				nextWanderingOrientation = envpkg.South
			} else {
				nextWanderingOrientation = envpkg.East
			}
		case envpkg.SouthWest:
			if chancesToChangeOrientation == 0 {
				nextWanderingOrientation = envpkg.South
			} else {
				nextWanderingOrientation = envpkg.West
			}
		}
	} else {
		nextWanderingOrientation = agt.orientation
	}
	// We go in the direction of the next orientation
	newObjective := agt.pos.Copy()
	for i := 0; i < agt.Speed; i++ {
		switch nextWanderingOrientation {
		case envpkg.North:
			newObjective.GoNorth(nil)
		case envpkg.South:
			newObjective.GoSouth(nil)
		case envpkg.East:
			newObjective.GoEast(nil)
		case envpkg.West:
			newObjective.GoWest(nil)
		case envpkg.NorthEast:
			newObjective.GoNorth(nil)
			newObjective.GoEast(nil)
		case envpkg.NorthWest:
			newObjective.GoNorth(nil)
			newObjective.GoWest(nil)
		case envpkg.SouthEast:
			newObjective.GoSouth(nil)
			newObjective.GoEast(nil)
		case envpkg.SouthWest:
			newObjective.GoSouth(nil)
			newObjective.GoWest(nil)
		}
	}
	return newObjective
}
