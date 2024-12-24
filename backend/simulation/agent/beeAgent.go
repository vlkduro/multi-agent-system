package agent

import (
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
	seenElems          []vision.SeenElem
	availablePositions []envpkg.Position
	objective          interface{}
}

type BeeAgentJson struct {
	ID           string             `json:"id"`
	Position     envpkg.Position    `json:"position"`
	Orientation  envpkg.Orientation `json:"orientation"`
	SeenElems    []vision.SeenElem  `json:"seenElems"`
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
	beeAgent.seenElems = []vision.SeenElem{}
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
		agt.objective = agt.hive
		return
	}
	var closestHornet *HornetAgent = nil
	for _, seen := range agt.seenElems {
		if seen.Elem != nil {
			switch elem := seen.Elem.(type) {
			case obj.Flower:
				if agt.objective == nil {
					agt.objective = elem
				}
			case HornetAgent:
				closestHornet = &elem
			}
		}
	}
	// Fleeing from the hornet
	if closestHornet != nil {
		distanceToHornet := agt.pos.DistanceFrom(*closestHornet.pos)
		for x := agt.pos.X - agt.Speed; x <= agt.pos.X+agt.Speed; x++ {
			for y := agt.pos.Y - agt.Speed; y <= agt.pos.Y+agt.Speed; y++ {
				if agt.env.IsValidPosition(x, y) {
					if agt.env.GetAt(x, y) == nil {
						pos := envpkg.Position{X: x, Y: y}
						if agt.pos.DistanceFrom(pos) > distanceToHornet {
							distanceToHornet = agt.pos.DistanceFrom(pos)
							agt.objective = pos
						}
					}
				}
			}
		}
	}
	if agt.objective == nil {
		isCloseToBorder := false
		var nextWanderingOrientation envpkg.Orientation
		// If we are close to the border, we go in the opposite direction
		for _, x := range []int{0, agt.env.GetMapDimension()} {
			distanceToBorder := int(agt.pos.DistanceFrom(envpkg.Position{X: x, Y: agt.pos.Y}))
			if distanceToBorder < agt.Speed {
				if x == 0 {
					nextWanderingOrientation = envpkg.East
				} else {
					nextWanderingOrientation = envpkg.West
				}
				isCloseToBorder = true
			}
		}
		if !isCloseToBorder {
			for _, y := range []int{0, agt.env.GetMapDimension()} {
				distanceToBorder := int(agt.pos.DistanceFrom(envpkg.Position{X: agt.pos.X, Y: y}))
				if distanceToBorder < agt.Speed {
					if y == 0 {
						nextWanderingOrientation = envpkg.South
					} else {
						nextWanderingOrientation = envpkg.North
					}
					isCloseToBorder = true
				}
			}
		}
		if !isCloseToBorder {
			// 1/4 chance to change orientation (2 possibilities for orientatino so we use 8)
			chancesToChangeOrientation := rand.Intn(8)
			if agt.orientation == envpkg.North || agt.orientation == envpkg.South {
				if chancesToChangeOrientation == 0 {
					agt.orientation = envpkg.East
				}
				if chancesToChangeOrientation == 1 {
					agt.orientation = envpkg.West
				}
			} else {
				if chancesToChangeOrientation == 0 {
					agt.orientation = envpkg.North
				}
				if chancesToChangeOrientation == 1 {
					agt.orientation = envpkg.South
				}
			}
			nextWanderingOrientation = agt.orientation
		}
		switch nextWanderingOrientation {
		case envpkg.North:
			agt.objective = envpkg.Position{X: agt.pos.X, Y: agt.pos.Y - agt.Speed}
		case envpkg.South:
			agt.objective = envpkg.Position{X: agt.pos.X, Y: agt.pos.Y + agt.Speed}
		case envpkg.East:
			agt.objective = envpkg.Position{X: agt.pos.X + agt.Speed, Y: agt.pos.Y}
		case envpkg.West:
			agt.objective = envpkg.Position{X: agt.pos.X - agt.Speed, Y: agt.pos.Y}
		}
	}
}

func (agt *BeeAgent) foragerAction() {
	if agt.objective != nil {
		switch obj := agt.objective.(type) {
		case envpkg.Position:
			if agt.pos.Equal(&obj) {
				agt.objective = nil
			} else {
				agt.gotoNextStepTowards(obj.Copy())
			}
		case obj.Flower:
			if agt.pos.Near(*obj.Position(), 1) {
				agt.nectar += obj.RetreiveNectar(agt.maxNectar - agt.nectar)
				agt.objective = nil
			} else {
				agt.gotoNextStepTowards(obj.Position().Copy())
			}
		case obj.Hive:
			if agt.pos.Near(*obj.Position(), 1) {
				obj.StoreNectar(agt.nectar)
				agt.nectar = 0
				agt.objective = nil
			} else {
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
			agt.pos.GoLeft(agt.env.GetMap())
			agt.orientation = envpkg.West
		} else {
			agt.pos.GoRight(agt.env.GetMap())
			agt.orientation = envpkg.East
		}
		if yFactor == 0 {
			agt.pos.GoUp(agt.env.GetMap())
			agt.orientation = envpkg.North
		} else {
			agt.pos.GoDown(agt.env.GetMap())
			agt.orientation = envpkg.South
		}
	}
}

func (agt *BeeAgent) ToJsonObj() interface{} {
	var objectivePos *envpkg.Position = nil
	if agt.objective != nil {
		switch obj := agt.objective.(type) {
		case envpkg.Position:
			objectivePos = &obj
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
