package environment

import (
	"sync"
)

type Environment struct {
	sync.RWMutex
	agts []Agent
	objs []Object
}

func NewEnvironment(agts []Agent, objs []Object) *Environment {
	return &Environment{agts: agts, objs: objs}
}

func (env *Environment) AddAgent(agt Agent) {
	env.agts = append(env.agts, agt)
}
