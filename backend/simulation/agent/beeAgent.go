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
	latestFlowerSeen   *obj.Flower
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

func NewBeeAgent(id string, env *envpkg.Environment, syncChan chan envpkg.AgentID, speed int, hive *obj.Hive, dob time.Time, maxnectar int, job job) *BeeAgent {
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
	beeAgent.pos = hive.Position().Copy()
	return beeAgent
}

func (agt *BeeAgent) Percept() {
	if agt.job == Forager {
		agt.foragerPerception()
	} else if agt.job == Worker {
		agt.workerPerception()
	}
}

func (agt *BeeAgent) Deliberate() {
	if agt.job == Forager {
		agt.foragerDeliberation()
	} else if agt.job == Worker {
		agt.workerDeliberation()
	}
}

func (agt *BeeAgent) Act() {
	if agt.job == Forager {
		agt.foragerAction()
	} else if agt.job == Worker {
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
	goesToHive := false
	if agt.nectar > 0 {
		decision := rand.Intn(101)
		chancesToGoToHive := agt.nectar / agt.maxNectar * 100
		if decision < chancesToGoToHive {
			fmt.Printf("[%s] Nectar full, going back to the hive\n", agt.id)
			agt.objective.TargetedElem = agt.hive
			agt.objective.Position = agt.hive.Position().Copy()
			agt.objective.Type = Hive
			goesToHive = true
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
					if elem.GetNectar() != 0 {
						if goesToHive {
							if agt.latestFlowerSeen == nil {
								fmt.Printf("[%s] Flower seen, preparing to tell hive !\n", agt.id)
								agt.latestFlowerSeen = elem
							}
							return
						} else {
							fmt.Printf("[%s] %v seen, going to it !\n", agt.id, elem.ID())
							agt.objective.TargetedElem = elem
							agt.objective.Position = elem.Position().Copy()
							agt.objective.Type = Flower
							hasAlreadySeenCloserFlower = true
						}
					} else {
						fmt.Printf("[%s] %v seen with no nectar, ignoring it !\n", agt.id, elem.ID())
					}
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
				fmt.Printf("[%s] Objectif reached %v\n", agt.id, objf.Position)
			} else {
				if _, ok := agt.env.GetAt(objf.Position.X, objf.Position.Y).(envpkg.IAgent); ok {
					// In some cases, the agent wants to go to a position where there is already an agent
					agt.objective.Type = None
					fmt.Printf("[%s] Agent at desired position, go wonder \n", agt.id)
				} else {
					agt.gotoNextStepTowards(objf.Position.Copy())
					if agt.pos.Equal(objf.Position) {
						agt.objective.Type = None
					}
					fmt.Printf("[%s] Go next position : %v \n", agt.id, objf.Position)
				}
			}
		case Flower:
			if flower, ok := objf.TargetedElem.(*obj.Flower); ok {
				if agt.pos.Near(objf.Position, 1) {
					agt.nectar += flower.RetreiveNectar(agt.maxNectar - agt.nectar)
					agt.objective.Type = None
					fmt.Printf("[%s] Nectar taken %v\n", agt.id, objf.TargetedElem.(envpkg.IObject).ID())
				} else {
					fmt.Printf("[%s] Going to flower : %v %v\n", agt.id, flower.ID(), flower.Position())
					agt.gotoNextStepTowards(flower.Position().Copy())
				}
			}
		case Hive:
			if agt.pos.Near(objf.Position, 1) {
				if hive, ok := objf.TargetedElem.(*obj.Hive); ok {
					if agt.nectar > 0 {
						hive.StoreNectar(agt.nectar)
						agt.nectar = 0
						agt.objective.Type = None
						fmt.Printf("[%s] Nectar stored %v\n", agt.id, objf.TargetedElem.(envpkg.IObject).ID())
					}
					foundFlower := agt.hive.GetLatestFlowerPos()
					// If the hive has a flower to tell the agent
					if foundFlower != nil {
						agt.objective.Position = foundFlower.Position().Copy()
						agt.objective.Type = Flower
						agt.objective.TargetedElem = foundFlower
						fmt.Printf("[%s] Hive has a flower to tell : %v\n", agt.id, foundFlower.ID())
					}
					// If the agent has seen a flower and is near the hive
					if agt.latestFlowerSeen != nil {
						fmt.Printf("[%s] I told Hive my flower : %v\n", agt.id, agt.latestFlowerSeen.ID())
						agt.hive.AddFlower(agt.latestFlowerSeen)
						agt.latestFlowerSeen = nil
					}
				}
			} else {
				fmt.Printf("[%s] Going to hive %v\n", agt.id, objf.TargetedElem.(envpkg.IObject).ID())
				agt.gotoNextStepTowards(objf.Position.Copy())
			}
		}
	}
}

func (agt *BeeAgent) workerPerception() {
	agt.visionFunc = vision.WorkerBeeVision
	agt.seenElems = agt.see()
}

func (agt *BeeAgent) workerDeliberation() {
	var closestHornet *HornetAgent = nil
	for _, seen := range agt.seenElems {
		if seen.Elem != nil {
			switch elem := (seen.Elem).(type) {
			case *HornetAgent:
				closestHornet = elem
			case *obj.Hive:
				beehivepos := agt.hive.Position()
				hiveseenpos := elem.Position()

				// check if bee is home
				if hiveseenpos.Equal(beehivepos) {
					// fmt.Printf("[%s] I work from home ! %v\n", agt.id, beehivepos)

					// check hive nectar quantity
					qNectarHive := agt.hive.GetNectar()
					if qNectarHive >= 4 {
						fmt.Printf("[%s] Objective set to ProduceHoney\n", agt.id)
						agt.objective.Type = ProduceHoney
					}
				}
				//default:
				//fmt.Printf("[%s] Unknown element seen : %v\n", agt.id, elem)
			}
		}
		if closestHornet != nil {
			break
		}
	}
}

// Template code to change, maybe think of inheritance for different
// bee types to avoid agmenting the complexity of the structure
// Perhaps inheriting a job component entirely could be a good idea
func (agt *BeeAgent) workerAction() {
	objf := &agt.objective
	fmt.Printf("[%s] Objective : %s\n", agt.id, objf.Type)
	// TODO: move to deliberating
	chancesToBecomeForager := rand.Intn(1000)
	if chancesToBecomeForager == 0 {
		agt.job = Forager
		xFactor := rand.Intn(2)
		yFactor := rand.Intn(2)
		agt.objective.Type = None
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
	} else {
		objf := &agt.objective
		if agt.objective.Type != None {
			switch typeObj := objf.Type; typeObj {
			case ProduceHoney:
				// fmt.Printf("[%s] Producing honey in hive\n", agt.id)
				/*if agt.hive.GetQNectar() > 4 {
					agt.hive.GetNectar(4)
					agt.hive.StoreHoney(1)
				}*/
			}
		}
		// produit 180mg de miel pour 600mg de nectar
		//  pour notre simulation :
		// 600mg nectar -> 150mg miel
		// 100mg nectar -> 25mg miel
		// 4mg nectar -> 1mg miel
		// 1 abeille 4mg nectar par tour
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
