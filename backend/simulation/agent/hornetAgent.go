package agent

import (
	"fmt"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/agent/vision"
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
	obj "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/object"
	"math/rand"
	"time"
)

// HornetAgent hérite de /simulation/agent/agent.go "struct Agent"
// Interface IAgent
type HornetAgent struct {
	Agent
	birthDate time.Time
	seenElems []*vision.SeenElem
	killCount int
}

type HornetAgentJson struct {
	ID          string             `json:"id"`
	Position    *envpkg.Position   `json:"position"`
	Orientation envpkg.Orientation `json:"orientation"`
	Objective   objective          `json:"objective"`
	SeenElems   []*vision.SeenElem `json:"seenElems"`
	KillCount   int                `json:"killCount"`
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
	hAgent.killCount = 0
	return hAgent
}

// A simple way to organize target priority for the hornet.
// If the hornet sees the Hive, and is surroundered by at least 2 of its kind, it will prioritize it.
// false : bee
// true : hive
func PriorityTarget(hornet HornetAgent) bool {
	nbHornet := 0
	if hornet.killCount >= 5 {
		return true
	}
	for _, seen := range hornet.seenElems {
		switch elem := seen.Elem.(type) {
		case *HornetAgent:
			if elem.ID() != hornet.ID() {
				fmt.Printf("HornetAgent %s seen by %s\n", elem.ID(), hornet.ID())
				nbHornet++
			}
		}
	}
	fmt.Printf("nbHornet = %d\n", nbHornet)
	return nbHornet >= 1
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

// The hornet agent always targets the closest bee unless :
// Sees a hive + is close to at least 2 of its kind
func (agt *HornetAgent) Deliberate() {
	distanceToTarget := float64(agt.env.GetMapDimension())
	if agt.objective.Type == Bee {
		if bee, ok := agt.objective.TargetedElem.(*BeeAgent); ok {
			distanceToTarget = bee.Position().DistanceFrom(agt.Position())
		}
	}
	if agt.objective.Type == Hive {
		if hive, ok := agt.objective.TargetedElem.(*obj.Hive); ok {
			distanceToTarget = hive.Position().DistanceFrom(agt.Position())
		}
	}
	for _, seen := range agt.seenElems {
		if seen.Elem != nil && seen.Elem != agt {
			priority := PriorityTarget(*agt)
			switch elem := seen.Elem.(type) {
			case *Agent:
				if elem.Type() == envpkg.Bee && !priority {
					fmt.Printf("[%s] Found a close bee (%s) \n", agt.id, elem.ID())
					distance := elem.Position().DistanceFrom(agt.Position())
					if distance < distanceToTarget {
						agt.objective.Type = Bee
						agt.objective.TargetedElem = elem
						agt.objective.Position = elem.Position().Copy()
						distanceToTarget = distance
					}
				}
			case *obj.Hive:
				if elem.TypeObject() == envpkg.Hive && priority {
					fmt.Printf("[%s] Found a close hive (%s) \n", agt.id, elem.ID())
					distance := elem.Position().DistanceFrom(agt.Position())
					if distance < distanceToTarget {
						agt.objective.Type = Hive
						agt.objective.TargetedElem = elem
						agt.objective.Position = elem.Position().Copy()
						distanceToTarget = distance

					}
				}
			}
		}
	}
	fmt.Printf(string(agt.objective.Type))
	if agt.objective.Type == None {
		agt.wander()
	}
}

func (agt *HornetAgent) Act() {
	objf := &agt.objective
	if objf.Type == None {
		return
	}
	switch objf.Type {
	case Position:
		if agt.pos.Equal(objf.Position) {
			agt.objective.Type = None
		} else if agt.pos.Equal(agt.lastPos) {
			// In some cases, the agent gets stuck on an unattaignable position
			agt.lastPos = objf.Position.Copy()
			agt.objective.Type = None
		} else {
			agt.gotoNextStepTowards(objf.Position.Copy())
			if agt.pos.Equal(objf.Position) {
				agt.objective.Type = None
			}
		}
	case Bee:
		bee := objf.TargetedElem.(*Agent)
		if bee != nil {
			fmt.Printf("[%s] Hornet attacking %s !\n", agt.id, bee.ID())
			agt.gotoNextStepTowards(bee.Position().Copy())
			if agt.Position().Near(bee.pos.Copy(), 1) {
				bee.Kill()
				agt.killCount++
				agt.objective.Type = None
				fmt.Printf("[%s] killed %s !!!, Has %d kill(s) \n", agt.id, bee.ID(), agt.killCount)
			}
		}

	case Hive:
		if agt.pos.Near(objf.Position, 1) {
			if hive, ok := objf.TargetedElem.(*obj.Hive); ok {
				if hive.IsAlive() {
					hive.Die()
				} else {
					fmt.Printf("Hive is already killed \n !")
				}
			}
		} else {
			fmt.Printf("Hornet [%s] going to kill hive %v\n", agt.id, objf.TargetedElem.(envpkg.IObject).ID())
			agt.gotoNextStepTowards(objf.Position.Copy())
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
