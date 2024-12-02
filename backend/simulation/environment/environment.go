package environment

import (
	"sync"
)

type Environment struct {
	sync.RWMutex
	agts []IAgent
	objs []IObject
}

func NewEnvironment(agts []IAgent, objs []IObject) *Environment {
	return &Environment{agts: agts, objs: objs}
}

func (env *Environment) AddAgent(agt IAgent) {
	env.agts = append(env.agts, agt)
}

func (env *Environment) AddObject(obj IObject) {
	env.objs = append(env.objs, obj)
}
