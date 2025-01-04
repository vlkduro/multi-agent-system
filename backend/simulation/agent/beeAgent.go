package agent

import (
	"fmt"
	"time"

	"math/rand"

	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/agent/vision"
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
	obj "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/object"
)

type job int

const (
	Worker = iota
	Guardian
	Forager
)

// BeeAgent hérite de /simulation/agent/agent.go "struct Agent"
// Interface IAgent
type BeeAgent struct {
	Agent
	hive               *obj.Hive
	birthDate          time.Time
	nectar             int
	maxNectar          int
	job                job
	seenElems          []*vision.SeenElem
	availablePositions []envpkg.Position
}

type BeeAgentJson struct {
	ID          string             `json:"id"`
	Position    *envpkg.Position   `json:"position"`
	Orientation envpkg.Orientation `json:"orientation"`
	SeenElems   []*vision.SeenElem `json:"seenElems"`
	MaxNectar   int                `json:"maxNectar"`
	Nectar      int                `json:"nectar"`
	Job         job                `json:"job"`
	Objective   objective          `json:"objective"`
}

func NewBeeAgent(id string, env *envpkg.Environment, syncChan chan bool, speed int, hive *obj.Hive, dob time.Time, maxnectar int, job job) *BeeAgent {
	beeAgent := &BeeAgent{}
	beeAgent.Agent = Agent{
		iagt:       beeAgent,
		id:         envpkg.AgentID(id),
		env:        env,
		syncChan:   syncChan,
		visionFunc: nil,
		speed:      speed,
		alive:      true,
		objective:  objective{Position: nil, Type: None},
	}
	beeAgent.hive = hive
	beeAgent.birthDate = dob
	beeAgent.maxNectar = maxnectar
	beeAgent.job = job
	beeAgent.nectar = 0
	beeAgent.availablePositions = []envpkg.Position{}
	beeAgent.seenElems = []*vision.SeenElem{}
	//beeAgent.pos = hive.Position().Copy()
	return beeAgent
}

func (agt *BeeAgent) Percept() {
	if agt.job == Forager {
		agt.foragerPerception()
	}
}

func (agt *BeeAgent) Deliberate() {
	if agt.job == Forager {
		agt.foragerDeliberation()
	}
}

func (agt *BeeAgent) Act() {
	if agt.job == Forager {
		agt.foragerAction()
	}
	if agt.job == Worker {
		agt.workerAction()
	}
}

func (agt *BeeAgent) hasFlowerObjective() bool {
	return agt.objective.Type == Flower
}

func (agt *BeeAgent) foragerPerception() {
	agt.visionFunc = vision.ExplorerBeeVision
	agt.seenElems = agt.see()
}

func (agt *BeeAgent) foragerDeliberation() {
	if agt.hasFlowerObjective() {
		return
	}
	if agt.nectar > 0 {
		decision := rand.Intn(101)
		chancesToGoToHive := agt.nectar / agt.maxNectar * 100
		if decision < chancesToGoToHive {
			fmt.Printf("[%s] Nectar full, going back to the hive\n", agt.id)
			agt.objective.TargetedElem = agt.hive
			agt.objective.Position = agt.hive.Position().Copy()
			agt.objective.Type = Hive
			return
		}
	}
	var closestHornet *HornetAgent = nil
	hasAlreadySeenCloserFlower := false
	for _, seen := range agt.seenElems {
		if seen.Elem != nil {
			switch elem := (seen.Elem).(type) {
			case *HornetAgent:
				fmt.Printf("[%s] Close to hornet %v\n", agt.id, elem.ID())
				closestHornet = elem
			case *obj.Flower:
				if !hasAlreadySeenCloserFlower {
					if elem.GetNectar() == 0 {
						fmt.Printf("[%s] Flower seen with no nectar, ignoring it !\n", agt.id)
						continue
					}
					fmt.Printf("[%s] Flower seen, going to it !\n", agt.id)
					agt.objective.TargetedElem = elem
					agt.objective.Position = elem.Position().Copy()
					agt.objective.Type = Flower
					hasAlreadySeenCloserFlower = true
				}
				//default:
				//fmt.Printf("[%s] Unknown element seen : %v\n", agt.id, elem)
			}
		}
		if closestHornet != nil {
			break
		}
	}
	// Fleeing from the hornet
	if closestHornet != nil {
		fmt.Printf("[%s] Hornet seen, fleeing in opposite direction\n", agt.id)
		agt.objective.Position = agt.pos.GetSymmetricOfPoint(*closestHornet.pos.Copy())
		agt.objective.Type = Position
		if agt.objective.Position.X < 0 {
			agt.objective.Position.X = 0
		} else if agt.objective.Position.X >= agt.env.GetMapDimension() {
			agt.objective.Position.X = agt.env.GetMapDimension() - 1
		}
		if agt.objective.Position.Y < 0 {
			agt.objective.Position.Y = 0
		} else if agt.objective.Position.Y >= agt.env.GetMapDimension() {
			agt.objective.Position.Y = agt.env.GetMapDimension() - 1
		}
	}
	// If has no objective, wander
	if agt.objective.Type == None {
		agt.wander()
	}
}

func (agt *BeeAgent) foragerAction() {
	objf := &agt.objective
	if agt.objective.Type != None {
		switch typeObj := objf.Type; typeObj {
		case Position:
			if agt.pos.Equal(objf.Position) {
				agt.objective.Type = None
			} else if agt.env.GetAt(objf.Position.X, objf.Position.Y) != nil && agt.pos.Near(objf.Position, 1) {
				// In some cases, the agent wants to go to a position where there is already an element (agent or object)
				agt.objective.Type = None
			} else {
				agt.gotoNextStepTowards(objf.Position.Copy())
				if agt.pos.Equal(objf.Position) {
					agt.objective.Type = None
				}
			}
		case Flower:
			if flower, ok := objf.TargetedElem.(*obj.Flower); ok {
				if agt.pos.Near(objf.Position, 1) {
					agt.nectar += flower.RetreiveNectar(agt.maxNectar - agt.nectar)
					agt.objective.Type = None
				} else {
					fmt.Printf("[%s] Going to flower %v\n", agt.id, objf.TargetedElem.(envpkg.IObject).ID())
					agt.gotoNextStepTowards(objf.Position.Copy())
				}
			}
		case Hive:
			if agt.pos.Near(objf.Position, 1) {
				if hive, ok := objf.TargetedElem.(*obj.Hive); ok {
					hive.StoreNectar(agt.nectar)
					agt.nectar = 0
					agt.objective.Type = None
				}
			} else {
				fmt.Printf("[%s] Going to hive %v\n", agt.id, objf.TargetedElem.(envpkg.IObject).ID())
				agt.gotoNextStepTowards(objf.Position.Copy())
			}
		}
	}
}

// Template code to change, maybe think of inheritance for different
// bee types to avoid agmenting the complexity of the structure
// Perhaps inheriting a job component entirely could be a good idea
func (agt *BeeAgent) workerAction() {
	chancesToBecomeForager := rand.Intn(10)
	if chancesToBecomeForager == 0 {
		agt.job = Forager
		xFactor := rand.Intn(2)
		yFactor := rand.Intn(2)

		agt.pos = agt.hive.Position().Copy()
		if xFactor == 0 {
			agt.goWest()
		} else {
			agt.goEast()
		}
		if yFactor == 0 {
			agt.goNorth()
		} else {
			agt.goSouth()
		}
	}
}

func (agt BeeAgent) Type() envpkg.AgentType {
	return envpkg.Bee
}

func (agt *BeeAgent) ToJsonObj() interface{} {
	return BeeAgentJson{ID: string(agt.id),
		Position:    agt.pos,
		Orientation: agt.orientation,
		SeenElems:   agt.seenElems,
		MaxNectar:   agt.maxNectar,
		Nectar:      agt.nectar,
		Job:         agt.job,
		Objective:   agt.objective,
	}
}
