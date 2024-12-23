package agent

import (
	"time"

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
	hive      obj.Hive
	birthDate time.Time
	nectar    int
	maxNectar int
	job       int
	seenElems []vision.SeenElem
	objective interface{}
}

type BeeAgentJson struct {
	ID        string `json:"id"`
	MaxNectar int    `json:"maxNectar"`
	Job       int    `json:"job"`
}

func NewBeeAgent(id string, env *envpkg.Environment, syncChan chan bool, speed int, r obj.Hive, dob time.Time, maxnectar int, job int) *BeeAgent {
	beeAgent := &BeeAgent{}
	beeAgent.Agent = Agent{
		iagt:       beeAgent,
		id:         envpkg.AgentID(id),
		env:        env,
		syncChan:   syncChan,
		visionFunc: nil,
		Speed:      speed,
	}
	beeAgent.hive = r
	beeAgent.birthDate = dob
	beeAgent.maxNectar = maxnectar
	beeAgent.job = job
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
}

func (agt *BeeAgent) ToJsonObj() interface{} {
	return BeeAgentJson{ID: string(agt.id), MaxNectar: agt.maxNectar, Job: agt.job}
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
}

func (agt *BeeAgent) foragerAction() {
	if agt.objective != nil {
		switch obj := agt.objective.(type) {
		case envpkg.Position:
			if agt.pos.Equal(&obj) {
				agt.objective = nil
			}
		case obj.Flower:
			if agt.pos.Near(*obj.Position(), 1) {
				agt.nectar += obj.RetreiveNectar(agt.maxNectar - agt.nectar)
				agt.objective = nil
			}
		}
	}
}
