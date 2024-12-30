package agent

import (
	"fmt"
	"math/rand"

	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/agent/vision"
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

type objectiveType string

const (
	None     objectiveType = "none"
	Position objectiveType = "position"
	Flower   objectiveType = "flower"
	Hive     objectiveType = "hive"
	Bee      objectiveType = "bee"
)

// BeeAgent structure to be marshalled in json
type objective struct {
	TargetedElem interface{}      `json:"targetedElem"`
	Position     *envpkg.Position `json:"position"`
	Type         objectiveType    `json:"type"`
}

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
	speed       int
	lastPos     *envpkg.Position
	alive       bool
	objective   objective
}

// Agent is launched as a microservice
func (agt *Agent) Start() {
	fmt.Printf("[%s] Starting Agent\n", agt.ID())
	go func() {
		for {
			run := <-agt.syncChan
			if !run || !agt.alive {
				agt.syncChan <- !run || !agt.alive
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

// Additionally we check the agent is on the position
func (agt *Agent) goNorth() bool {
	success := agt.pos.GoNorth(agt.env.GetMap(), agt)
	agt.orientation = envpkg.North
	if agt.env.GetAt(agt.pos.X, agt.pos.Y) == nil {
		agt.env.GetMap()[agt.pos.X][agt.pos.Y] = agt
	}
	return success
}

func (agt *Agent) forceAgentInPlace(env *envpkg.Environment) {
	if env.GetAt(agt.pos.X, agt.pos.Y) == nil {
		env.GetMap()[agt.pos.X][agt.pos.Y] = agt
	}
}

func (agt *Agent) goSouth() bool {
	success := agt.pos.GoSouth(agt.env.GetMap(), agt)
	agt.orientation = envpkg.South
	agt.forceAgentInPlace(agt.env)
	return success
}

func (agt *Agent) goEast() bool {
	success := agt.pos.GoEast(agt.env.GetMap(), agt)
	agt.orientation = envpkg.East
	agt.forceAgentInPlace(agt.env)
	return success
}

func (agt *Agent) goWest() bool {
	success := agt.pos.GoWest(agt.env.GetMap(), agt)
	agt.orientation = envpkg.West
	agt.forceAgentInPlace(agt.env)
	return success
}

func (agt *Agent) goNorthEast() bool {
	success := agt.pos.GoNorthEast(agt.env.GetMap(), agt)
	agt.orientation = envpkg.NorthEast
	agt.forceAgentInPlace(agt.env)
	return success
}

func (agt *Agent) goSouthEast() bool {
	success := agt.pos.GoSouthEast(agt.env.GetMap(), agt)
	agt.orientation = envpkg.SouthEast
	agt.forceAgentInPlace(agt.env)
	return success
}

func (agt *Agent) goNorthWest() bool {
	success := agt.pos.GoNorthWest(agt.env.GetMap(), agt)
	agt.orientation = envpkg.NorthWest
	agt.forceAgentInPlace(agt.env)
	return success
}

func (agt *Agent) goSouthWest() bool {
	success := agt.pos.GoSouthWest(agt.env.GetMap(), agt)
	agt.orientation = envpkg.SouthWest
	agt.forceAgentInPlace(agt.env)
	return success
}

// https://web.archive.org/web/20171022224528/http://www.policyalmanac.org:80/games/aStarTutorial.htm
func (agt *Agent) gotoNextStepTowards(pos *envpkg.Position) {
	chain := agt.env.PathFinding(agt.pos, pos, agt.speed)
	// We remove the first element who is the current position of the agent
	if len(chain) > 0 && chain[0].Equal(agt.pos) {
		chain = chain[1:]
	}
	//fmt.Printf("\n[%s] Going to [%d %d] : [%d %d] -> ", agt.id, pos.X, pos.Y, agt.pos.X, agt.pos.Y)
	for i := 0; i < agt.speed && i < len(chain); i++ {
		pos := chain[i]
		if agt.pos.X < pos.X {
			if agt.pos.Y < pos.Y {
				agt.goSouthEast()
			} else if agt.pos.Y > pos.Y {
				agt.goNorthEast()
			} else {
				agt.goEast()
			}
		} else if agt.pos.X > pos.X {
			if agt.pos.Y < pos.Y {
				agt.goSouthWest()
			} else if agt.pos.Y > pos.Y {
				agt.goNorthWest()
			} else {
				agt.goWest()
			}
		} else {
			if agt.pos.Y < pos.Y {
				agt.goSouth()
			} else if agt.pos.Y > pos.Y {
				agt.goNorth()
			}
		}
	}
}

func (agt *Agent) wander() {
	if agt.pos == nil {
		return
	}
	var closestBorder *envpkg.Position = nil
	minBorderDistance := agt.pos.DistanceFrom(&envpkg.Position{X: 0, Y: 0})
	isCloseToBorder := false
	// If we are close to the border, we go in the opposite direction
	// We put -1 in the list to test the flat border cases
	for _, x := range []int{-1, 0, agt.env.GetMapDimension() - 1} {
		for _, y := range []int{-1, 0, agt.env.GetMapDimension() - 1} {
			isCorner := true
			if x == -1 && y == -1 {
				continue
			}
			// Flat north or south border
			if x == -1 {
				x = agt.pos.X
				isCorner = false
			}
			if y == -1 {
				// Flat west or east border
				y = agt.pos.Y
				isCorner = false
			}
			// Corner case
			distance := agt.pos.DistanceFrom(&envpkg.Position{X: x, Y: y})
			isCloseToBorder = distance < float64(agt.speed)
			// We allow some leeway to border cases to avoid getting stuck
			if (isCloseToBorder && distance < minBorderDistance) || (isCloseToBorder && isCorner && distance <= minBorderDistance) {
				closestBorder = &envpkg.Position{X: x, Y: y}
				minBorderDistance = distance
			}
		}
	}
	if closestBorder != nil {
		// If we are too close to the border, we go to the opposite side
		keepAwayFromBorderPos := agt.pos.Copy()
		if closestBorder.X == 0 {
			keepAwayFromBorderPos.GoEast(nil, nil)
		} else if closestBorder.X == agt.env.GetMapDimension()-1 {
			keepAwayFromBorderPos.GoWest(nil, nil)
		}
		if closestBorder.Y == 0 {
			keepAwayFromBorderPos.GoSouth(nil, nil)
		} else if closestBorder.Y == agt.env.GetMapDimension()-1 {
			keepAwayFromBorderPos.GoNorth(nil, nil)
		}
		agt.objective.Position = keepAwayFromBorderPos.Copy()
		agt.objective.Type = Position
		fmt.Printf("[%s] Too close to border (%d %d), going to (%d %d)\n", agt.id, closestBorder.X, closestBorder.Y, agt.objective.Position.X, agt.objective.Position.Y)
	} else {
		// While we don't have an objective, we wander
		for agt.objective.Type == None {
			newObjective := agt.getNextWanderingPosition()
			elemAtObjective := agt.env.GetAt(newObjective.X, newObjective.Y)
			if elemAtObjective != nil {
				continue
			}
			fmt.Printf("[%s] Wandering towards %v\n", agt.id, *newObjective)
			agt.objective.Type = Position
			agt.objective.Position = newObjective.Copy()
		}
	}
}

func (agt *Agent) getNextWanderingPosition() *envpkg.Position {
	surroundings := agt.pos.GetNeighbours(agt.speed)
	// We remove the positions that are occupied
	removeCpt := 0
	for i := 0; i < len(surroundings); i++ {
		idx := i - removeCpt
		if agt.env.GetAt(surroundings[idx].X, surroundings[idx].Y) != nil || surroundings[idx].Equal(agt.lastPos) {
			surroundings = append(surroundings[:idx], surroundings[idx+1:]...)
			removeCpt++
		}
	}
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
	for i := 0; i < agt.speed; i++ {
		switch nextWanderingOrientation {
		case envpkg.North:
			newObjective.GoNorth(nil, nil)
		case envpkg.South:
			newObjective.GoSouth(nil, nil)
		case envpkg.East:
			newObjective.GoEast(nil, nil)
		case envpkg.West:
			newObjective.GoWest(nil, nil)
		case envpkg.NorthEast:
			newObjective.GoNorthEast(nil, nil)
		case envpkg.NorthWest:
			newObjective.GoNorthWest(nil, nil)
		case envpkg.SouthEast:
			newObjective.GoSouthEast(nil, nil)
		case envpkg.SouthWest:
			newObjective.GoSouthWest(nil, nil)
		}
	}
	// Find the closest available position in surroundings
	closestPosition := newObjective.Copy()
	minDistance := agt.pos.DistanceFrom(newObjective)
	for _, pos := range surroundings {
		distance := pos.DistanceFrom(newObjective)
		if distance < minDistance {
			closestPosition = pos.Copy()
			minDistance = distance
		}
	}
	newObjective = closestPosition

	// We add the new last position to avoid cycles
	agt.lastPos = newObjective.Copy()

	return newObjective
}

func (agt *Agent) die() {
	agt.alive = false
	agt.pos = nil
}
