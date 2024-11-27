package simulation

import (
	env "./environment"
)

type Simulation struct {
	env     *env.Environment
	agents  []env.IAgent
	objects []env.IObject
}

func (simu *Simulation) Run()
func (simu *Simulation) log()
func (simu *Simulation) print()
