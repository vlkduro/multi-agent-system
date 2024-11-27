package environment

import (
	"sync"
)

type Environment struct {
	sync.RWMutex
	agts    []Agent
	objects []Object
}
