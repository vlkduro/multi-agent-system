package environment

import (
	"sync"

	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

var mapDimension int

type Environment struct {
	sync.Mutex
	agts []IAgent
	objs []IObject
	grid [][]interface{}
}

type EnvironmentJson struct {
	MapDimension int             `json:"mapDimension"`
	Grid         [][]interface{} `json:"grid"`
}

func NewEnvironment(agts []IAgent, objs []IObject) *Environment {
	mapDimension = utils.GetMapDimension()
	grid := make([][]interface{}, mapDimension)
	for i := range grid {
		grid[i] = make([]interface{}, mapDimension)
	}
	return &Environment{agts: agts, objs: objs, grid: grid}
}

func (env *Environment) GetAt(x int, y int) interface{} {
	if x < 0 || y < 0 || x >= mapDimension || y >= mapDimension {
		return nil
	}
	return env.grid[x][y]
}

func (env *Environment) GetMap() [][]interface{} {
	return env.grid
}

func (env *Environment) AddAgent(agt IAgent) bool {
	pos := agt.Position()
	if env.GetAt(pos.X, pos.Y) != nil {
		return false
	}
	env.agts = append(env.agts, agt)
	env.grid[pos.X][pos.Y] = agt
	return true
}

func (env *Environment) AddObject(obj IObject) bool {
	pos := obj.Position()
	if env.GetAt(pos.X, pos.Y) != nil {
		return false
	}
	env.objs = append(env.objs, obj)
	env.grid[pos.X][pos.Y] = obj
	return true
}

func (env *Environment) ToJsonObj() interface{} {
	grid := make([][]interface{}, mapDimension)
	for i := range grid {
		grid[i] = make([]interface{}, mapDimension)
	}

	for x := range env.grid {
		for y := range env.grid[x] {
			gridObj := env.grid[x][y]
			if gridObj != nil {
				switch obj := gridObj.(type) {
				case IAgent:
					grid[x][y] = obj.ToJsonObj()
				case IObject:
					grid[x][y] = obj.ToJsonObj()
				}
			}
		}
	}

	return EnvironmentJson{MapDimension: mapDimension, Grid: grid}
}
