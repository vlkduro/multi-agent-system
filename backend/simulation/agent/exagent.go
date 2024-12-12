package agent

import (
	"fmt"
	"math/rand/v2"

	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/agent/vision"
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
)

const VisionRange = 6

// ExAgent est un exemple d'agent
type ExAgent struct {
	Agent
	value    int
	toAdd    int
	movement envpkg.Orientation
	vision   string
	// for debug
	seenElems []vision.SeenElem
}

type ExAgentJson struct {
	ID          string             `json:"id"`
	Value       int                `json:"value"`
	ToAdd       int                `json:"toAdd"`
	Movement    envpkg.Orientation `json:"movement"`
	Position    envpkg.Position    `json:"position"`
	Orientation envpkg.Orientation `json:"orientation"`
	Vision      string             `json:"vision"`
	SeenElems   []vision.SeenElem  `json:"seenElems"`
}

func NewExAgent(id string, pos *envpkg.Position, env *envpkg.Environment, syncChan chan bool) *ExAgent {
	exagt := &ExAgent{}
	exagt.Agent = Agent{
		iagt:        exagt,
		id:          envpkg.AgentID(id),
		pos:         pos.Copy(),
		orientation: envpkg.East,
		env:         env,
		visionFunc:  vision.ExplorerBeeVision,
		syncChan:    syncChan,
	}

	return exagt
}

func (agt *ExAgent) Percept() {
	agt.vision = ""
	logPos := func(item string, x int, y int) {
		agt.vision += fmt.Sprintf("%s[%d;%d] ", item, x, y)
	}
	logPos("ME", agt.pos.X, agt.pos.Y)
	agt.seenElems = agt.see()
	for _, seen := range agt.seenElems {
		if seen.Elem != nil {
			switch v := seen.Elem.(type) {
			case envpkg.IObject:
				logPos(string(v.ID()), seen.Pos.X, seen.Pos.Y)
			case envpkg.IAgent:
				logPos(string(v.ID()), seen.Pos.X, seen.Pos.Y)
			}
		} else {
			logPos(".", seen.Pos.X, seen.Pos.Y)
		}
	}

}

func (agt *ExAgent) Deliberate() {
	factor := 1
	if rand.IntN(2) == 0 {
		factor = -1
	}
	agt.toAdd = factor * rand.IntN(100)
	rnd := rand.IntN(2)
	// Stupid movement deliberation
	if agt.toAdd < 0 {
		if rnd == 0 {
			agt.movement = envpkg.North
		} else {
			agt.movement = envpkg.South
		}
	} else {
		if rnd == 0 {
			agt.movement = envpkg.West
		} else {
			agt.movement = envpkg.East
		}
	}
}

func (agt *ExAgent) Act() {
	agt.value += agt.toAdd
	switch agt.movement {
	case envpkg.North:
		agt.pos.GoUp(agt.env.GetMap())
	case envpkg.East:
		agt.pos.GoRight(agt.env.GetMap())
	case envpkg.South:
		agt.pos.GoDown(agt.env.GetMap())
	case envpkg.West:
		agt.pos.GoLeft(agt.env.GetMap())
	}
	agt.orientation = agt.movement
}

func (agt ExAgent) ToJsonObj() interface{} {
	return ExAgentJson{ID: string(agt.id),
		Value:       agt.value,
		ToAdd:       agt.toAdd,
		Movement:    agt.movement,
		Position:    *agt.pos,
		Orientation: agt.orientation,
		Vision:      agt.vision,
		SeenElems:   agt.seenElems}
}
