package environment

import (
	"sync"
)

type Environment struct {
	sync.RWMutex
	agts    []IAgent
	objects []IObject
}
