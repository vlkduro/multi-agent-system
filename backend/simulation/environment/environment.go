package environment

import (
	"sync"
)

type Environment struct {
	sync.RWMutex
	agts []IAgent
	objs []Object
}

func NewEnvironment(agts []IAgent, objs []Object) *Environment {
	return &Environment{agts: agts, objs: objs}
}

func (env *Environment) AddAgent(agt IAgent) {
	env.agts = append(env.agts, agt)
}
