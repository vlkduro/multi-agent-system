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
	// If simple Position, make it a pointer
	objective interface{}
}

type BeeAgentJson struct {
	ID           string             `json:"id"`
	Position     envpkg.Position    `json:"position"`
	Orientation  envpkg.Orientation `json:"orientation"`
	SeenElems    []*vision.SeenElem `json:"seenElems"`
	MaxNectar    int                `json:"maxNectar"`
	Nectar       int                `json:"nectar"`
	Job          job                `json:"job"`
	ObjectivePos *envpkg.Position   `json:"objectivePos"`
}

func NewBeeAgent(id string, env *envpkg.Environment, syncChan chan bool, speed int, hive *obj.Hive, dob time.Time, maxnectar int, job job) *BeeAgent {
	beeAgent := &BeeAgent{}
	beeAgent.Agent = Agent{
		iagt:       beeAgent,
		id:         envpkg.AgentID(id),
		env:        env,
		syncChan:   syncChan,
		visionFunc: nil,
		Speed:      speed,
	}
	beeAgent.hive = hive
	beeAgent.birthDate = dob
	beeAgent.maxNectar = maxnectar
	beeAgent.job = job
	beeAgent.nectar = 0
	beeAgent.objective = nil
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

func (agt *BeeAgent) foragerPerception() {
	agt.visionFunc = vision.ExplorerBeeVision
	agt.seenElems = agt.see()
}

func (agt *BeeAgent) foragerDeliberation() {
	if agt.nectar == agt.maxNectar {
		fmt.Printf("[%s] Nectar full, going back to the hive\n", agt.id)
		agt.objective = agt.hive
		return
	}
	var closestHornet *HornetAgent = nil
	hasAlreadySeenCloserFlower := false
	for _, seen := range agt.seenElems {
		if seen.Elem != nil {
			switch elem := (seen.Elem).(type) {
			case HornetAgent:
				closestHornet = &elem
			case obj.Flower:
				if !hasAlreadySeenCloserFlower {
					fmt.Printf("[%s] Flower seen, going to it !\n", agt.id)
					agt.objective = elem
					hasAlreadySeenCloserFlower = true
				}
			}
		}
		if closestHornet != nil {
			break
		}
	}
	// Fleeing from the hornet
	if closestHornet != nil {
		fmt.Printf("[%s] Hornet seen, fleeing in opposite direction\n", agt.id)
		agt.objective = agt.pos.GetSymmetricOfPoint(*closestHornet.pos.Copy())
	}
	// If has no objective, wander
	if agt.objective == nil {
		var closestBorder *envpkg.Position = nil
		minBorderDistance := agt.pos.DistanceFrom(envpkg.Position{X: 0, Y: 0})
		isCloseToBorder := false
		var nextWanderingOrientation envpkg.Orientation
		// If we are close to the border, we go in the opposite direction
		for _, x := range []int{0, agt.env.GetMapDimension()} {
			for _, y := range []int{0, agt.env.GetMapDimension()} {
				distance := agt.pos.DistanceFrom(envpkg.Position{X: x, Y: y})
				isCloseToBorder = distance < float64(agt.Speed)
				if isCloseToBorder && distance <= minBorderDistance {
					closestBorder = &envpkg.Position{X: x, Y: y}
					minBorderDistance = distance
				}
			}
		}
		if closestBorder != nil {
			fmt.Printf("[%s] Too close to border\n", agt.id)
			agt.objective = agt.pos.GetSymmetricOfPoint(*closestBorder)
		} else {
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
			fmt.Printf("[%s] Wandering towards %v\n", agt.id, *newObjective)
			agt.objective = newObjective.Copy()
		}
	}
}

func (agt *BeeAgent) foragerAction() {
	if agt.objective != nil {
		switch obj := agt.objective.(type) {
		case *envpkg.Position:
			if agt.pos.Equal(obj) {
				agt.objective = nil
			} else {
				fmt.Printf("[%s] Going to %v\n", agt.id, *obj)
				agt.gotoNextStepTowards(obj.Copy())
			}
		case obj.Flower:
			if agt.pos.Near(*obj.Position(), 1) {
				agt.nectar += obj.RetreiveNectar(agt.maxNectar - agt.nectar)
				agt.objective = nil
			} else {
				fmt.Printf("[%s] Going to flower %v\n", agt.id, obj.ID())
				agt.gotoNextStepTowards(obj.Position().Copy())
			}
		case obj.Hive:
			if agt.pos.Near(*obj.Position(), 1) {
				obj.StoreNectar(agt.nectar)
				agt.nectar = 0
				agt.objective = nil
			} else {
				fmt.Printf("[%s] Going to hive %v\n", agt.id, obj.ID())
				agt.gotoNextStepTowards(obj.Position().Copy())
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
	var objectivePos *envpkg.Position = nil
	if agt.objective != nil {
		switch obj := agt.objective.(type) {
		case *envpkg.Position:
			objectivePos = obj
		case obj.Flower:
			objectivePos = obj.Position().Copy()
		case obj.Hive:
			objectivePos = obj.Position().Copy()
		}
	}
	return BeeAgentJson{ID: string(agt.id),
		Position:     *agt.pos,
		Orientation:  agt.orientation,
		SeenElems:    agt.seenElems,
		MaxNectar:    agt.maxNectar,
		Nectar:       agt.nectar,
		Job:          agt.job,
		ObjectivePos: objectivePos,
	}
}
