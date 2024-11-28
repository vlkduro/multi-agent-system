package simulation

import (
	"fmt"
	"log"
	"time"

	agt "./agent"
	envpkg "./environment"
)

type Simulation struct {
	env  *envpkg.Environment
	agts []envpkg.Agent
	objs []envpkg.Object
}

func NewSimulation(nagt int, nobj int) *Simulation {
	simu := &Simulation{}
	env := envpkg.NewEnvironment([]envpkg.Agent{}, []envpkg.Object{})

	simu.env = env

	for i := 0; i < nagt; i++ {
		// création de l'agent
		id := fmt.Sprintf("Agent #%d", i)
		syncChan := make(chan bool)
		agt := agt.NewExAgent(id, simu.env, syncChan)

		// ajout de l'agent à la simulation
		simu.agts = append(simu.agts, agt)

		// ajout de l'agent à l'environnement
		simu.env.AddAgent(agt)
	}
	return simu
}

func (simu *Simulation) Run() {
	log.Printf("Démarrage de la simulation")

	// Démarrage du micro-service de Log
	go simu.log()
	// Démarrage du micro-service d'affichage
	go simu.print()

	// Démarrage des agents
	for i := range simu.agts {
		simu.agts[i].Start()
	}

	// Boucle de simulation
	for _, agt := range simu.agts {
		go func(agt envpkg.Agent) {
			for {
				c := agt.GetSyncChan()
				c <- false
				time.Sleep(1 * time.Millisecond) // attente avant de relancer l'agent
				<-c
			}
		}(agt)
	}

	time.Sleep(time.Second / 120) // 120 tps

	log.Printf("Fin de la simulation")
}

// Intention d'en faire un microservice
func (simu *Simulation) log()

// Intention d'en faire un microservice
func (simu *Simulation) print() {
	for {
		fmt.Printf("TODO")
		time.Sleep(time.Second / 60) // 60 fps
	}

}
