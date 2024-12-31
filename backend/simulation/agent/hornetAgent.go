package agent

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/agent/vision"
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

// HornetAgent hérite de /simulation/agent/agent.go "struct Agent"
// Interface IAgent
type HornetAgent struct {
	Agent
	birthDate time.Time
	seenElems []*vision.SeenElem
}

type HornetAgentJson struct {
	ID          string             `json:"id"`
	Position    *envpkg.Position   `json:"position"`
	Orientation envpkg.Orientation `json:"orientation"`
	Objective   objective          `json:"objective"`
	SeenElems   []*vision.SeenElem `json:"seenElems"`
}

func NewHornetAgent(id string, env *envpkg.Environment, syncChan chan bool, s int) *HornetAgent {
	hAgent := &HornetAgent{}
	hAgent.Agent = Agent{
		iagt:       hAgent,
		id:         envpkg.AgentID(id),
		env:        env,
		syncChan:   syncChan,
		speed:      s,
		visionFunc: vision.HornetVision,
		alive:      true,
		objective:  objective{Position: nil, Type: None},
	}
	hAgent.birthDate = time.Now()
	return hAgent
}

func (agt *HornetAgent) Percept() {
	if agt.pos == nil {
		chance := rand.Intn(2)
		if chance == 0 {
			// either spawns at the top or the bottom or the left or the right
			chance = rand.Intn(4)
			x := 0
			y := 0
			switch chance {
			case 0:
				x = 0
				y = rand.Intn(agt.env.GetMapDimension())
				agt.orientation = envpkg.East
			case 1:
				x = agt.env.GetMapDimension() - 1
				y = rand.Intn(agt.env.GetMapDimension())
				agt.orientation = envpkg.West
			case 2:
				x = rand.Intn(agt.env.GetMapDimension())
				y = 0
				agt.orientation = envpkg.South
			case 3:
				x = rand.Intn(agt.env.GetMapDimension())
				y = agt.env.GetMapDimension() - 1
				agt.orientation = envpkg.North
			}
			agt.pos = envpkg.NewPosition(x, y, agt.env.GetMapDimension(), agt.env.GetMapDimension())
			agt.env.AddAgent(agt)
			fmt.Printf("[%s] Hornet spawned at %d %d\n", agt.id, x, y)
			agt.objective.Type = None
		}
	}
	if agt.pos != nil {
		agt.seenElems = agt.see()
	}
}

// The hornet agent always targets the closest bee
func (agt *HornetAgent) Deliberate() {
	distanceToTarget := float64(agt.env.GetMapDimension())
	if agt.objective.Type == Bee {
		if bee, ok := agt.objective.TargetedElem.(*BeeAgent); ok {
			distanceToTarget = bee.Position().DistanceFrom(agt.Position())
		}
	}
	for _, seen := range agt.seenElems {
		if seen.Elem != nil && seen.Elem != agt {
			switch elem := seen.Elem.(type) {
			case *Agent:
				if elem.Type() == envpkg.Bee {
					fmt.Printf("[%s] Found a close bee (%s) \n", agt.id, elem.ID())
					distance := elem.Position().DistanceFrom(agt.Position())
					if distance < distanceToTarget {
						agt.objective.Type = Bee
						agt.objective.TargetedElem = elem
						agt.objective.Position = elem.Position().Copy()
						distanceToTarget = distance
					}
				}
			}
		}
	}
	if agt.objective.Type == None {
		agt.wander()
	}
}

func (agt *HornetAgent) Act() {
	if agt.objective.Type == None {
		return
	}
	switch agt.objective.Type {
	case Position:
		agt.gotoNextStepTowards(agt.objective.Position.Copy())
		if agt.Position().Equal(agt.objective.Position) {
			agt.objective.Type = None
		}
	case Bee:
		bee := agt.objective.TargetedElem.(*Agent)
		fmt.Printf("[%s] Hornet attacking %s !\n", agt.id, bee.ID())
		agt.gotoNextStepTowards(bee.Position().Copy())
		if agt.Position().Near(bee.pos.Copy(), 1) {
			bee.Kill()
			agt.objective.Type = None
			fmt.Printf("[%s] Hornet killed %s !!!\n", agt.id, bee.ID())
		}
	}
}

func (agt HornetAgent) Type() envpkg.AgentType {
	return envpkg.Hornet
}

func (agt *HornetAgent) ToJsonObj() interface{} {
	return HornetAgentJson{
		ID:          string(agt.id),
		Position:    agt.pos,
		Orientation: agt.orientation,
		Objective:   agt.objective,
		SeenElems:   agt.seenElems,
	}
}
