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
	objective          objective
}

type objectiveType string

const (
	None     objectiveType = "none"
	Position objectiveType = "position"
	Flower   objectiveType = "flower"
	Hive     objectiveType = "hive"
)

// BeeAgent structure to be marshalled in json
type objective struct {
	TargetedElem envpkg.IObject
	Position     *envpkg.Position `json:"position"`
	Type         objectiveType    `json:"type"`
}

type BeeAgentJson struct {
	ID          string             `json:"id"`
	Position    envpkg.Position    `json:"position"`
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
	}
	beeAgent.hive = hive
	beeAgent.birthDate = dob
	beeAgent.maxNectar = maxnectar
	beeAgent.job = job
	beeAgent.nectar = 0
	beeAgent.objective = objective{Position: nil, Type: None}
	beeAgent.availablePositions = []envpkg.Position{}
	beeAgent.seenElems = []*vision.SeenElem{}
	beeAgent.pos = hive.Position().Copy()
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
	if agt.objective.Type == Flower {
		return true
	}
	return false
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
			default:
				fmt.Printf("[%s] Unknown element seen : %v\n", agt.id, elem)
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
	}
	// If has no objective, wander
	if agt.objective.Type == None {
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
				keepAwayFromBorderPos.GoEast(nil)
			} else if closestBorder.X == agt.env.GetMapDimension()-1 {
				keepAwayFromBorderPos.GoWest(nil)
			}
			if closestBorder.Y == 0 {
				keepAwayFromBorderPos.GoSouth(nil)
			} else if closestBorder.Y == agt.env.GetMapDimension()-1 {
				keepAwayFromBorderPos.GoNorth(nil)
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
}

func (agt *BeeAgent) foragerAction() {
	objf := &agt.objective
	if agt.objective.Type != None {
		switch typeObj := objf.Type; typeObj {
		case Position:
			if agt.pos.Equal(objf.Position) {
				agt.objective.Type = None
			} else {
				agt.gotoNextStepTowards(objf.Position.Copy())
			}
		case Flower:
			if flower, ok := objf.TargetedElem.(*obj.Flower); ok {
				if agt.pos.Near(objf.Position, 1) {
					agt.nectar += flower.RetreiveNectar(agt.maxNectar - agt.nectar)
					agt.objective.Type = None
				} else {
					fmt.Printf("[%s] Going to flower %v\n", agt.id, objf.TargetedElem.ID())
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
				fmt.Printf("[%s] Going to hive %v\n", agt.id, objf.TargetedElem.ID())
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
			agt.pos.GoWest(agt.env.GetMap())
			agt.orientation = envpkg.West
		} else {
			agt.pos.GoEast(agt.env.GetMap())
			agt.orientation = envpkg.East
		}
		if yFactor == 0 {
			agt.pos.GoNorth(agt.env.GetMap())
			agt.orientation = envpkg.North
		} else {
			agt.pos.GoSouth(agt.env.GetMap())
			agt.orientation = envpkg.South
		}
	}
}

func (agt *BeeAgent) ToJsonObj() interface{} {
	return BeeAgentJson{ID: string(agt.id),
		Position:    *agt.pos,
		Orientation: agt.orientation,
		SeenElems:   agt.seenElems,
		MaxNectar:   agt.maxNectar,
		Nectar:      agt.nectar,
		Job:         agt.job,
		Objective:   agt.objective,
	}
}
