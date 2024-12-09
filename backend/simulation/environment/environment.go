package environment

import (
	"sync"
)

type Environment struct {
	sync.Mutex
	agts []IAgent
	objs []IObject
	grid [][]interface{}
}

func NewEnvironment(agts []IAgent, objs []IObject) *Environment {
	grid := make([][]interface{}, 12)
	for i := range grid {
		grid[i] = make([]interface{}, 12)
	}
	return &Environment{agts: agts, objs: objs, grid: grid}
}

func (env *Environment) GetAt(x int, y int) interface{} {
	if x < 0 || y < 0 {
		return nil
	}
	return env.grid[x][y]
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
