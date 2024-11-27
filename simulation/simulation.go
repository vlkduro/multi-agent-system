package simulation

import (
	env "./environment"
)

type Simulation struct {
	env     *env.Environment
	agents  []env.Agent
	objects []env.Object
}

func (simu *Simulation) Run()
func (simu *Simulation) log()
func (simu *Simulation) print()
