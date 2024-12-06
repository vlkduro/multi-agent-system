package agent

import (
	"fmt"
	"math/rand/v2"

	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

const VisionRange = 6

// ExAgent est un exemple d'agent
type ExAgent struct {
	Agent
	value    int
	decision int
	vision   string
}

type ExAgentJson struct {
	ID       string          `json:"id"`
	Value    int             `json:"value"`
	Decision int             `json:"decision"`
	Position envpkg.Position `json:"position"`
	Vision   string          `json:"vision"`
}

func NewExAgent(id string, pos *envpkg.Position, env *envpkg.Environment, syncChan chan bool) *ExAgent {
	exagt := &ExAgent{}
	exagt.Agent = Agent{
		iagt:     exagt,
		id:       envpkg.AgentID(id),
		pos:      pos.Copy(),
		env:      env,
		syncChan: syncChan,
	}

	return exagt
}

func (agt *ExAgent) Percept() {
	agt.vision = ""
	for x := -int(VisionRange / 2); x < int(VisionRange/2); x++ {
		for y := -int(VisionRange / 2); y < int(VisionRange/2); y++ {
			if x == 0 && y == 0 {
				agt.vision += "ME"
				continue
			}
			if seen := agt.env.GetAt(agt.pos.X+x, agt.pos.Y+y); seen != nil {
				switch v := seen.(type) {
				case envpkg.IObject:
					agt.vision += fmt.Sprintf("%s[%d;%d] ", string(v.ID()), v.Position().X, v.Position().Y)
				case envpkg.IAgent:
					agt.vision += fmt.Sprintf("%s[%d;%d] ", string(v.ID()), v.Position().X, v.Position().Y)
				}
			} else {
				agt.vision += "."
			}
		}
		agt.vision += "\n"
	}
}

func (agt *ExAgent) Deliberate() {
	factor := 1
	if rand.IntN(2) == 0 {
		factor = -1
	}
	agt.decision = factor * rand.IntN(100)
}

func (agt *ExAgent) Act() {
	agt.value += agt.decision
}

func (agt ExAgent) ToJsonObj() interface{} {
	return ExAgentJson{ID: string(agt.id), Value: agt.value, Decision: agt.decision, Position: *agt.pos, Vision: agt.vision}
}
