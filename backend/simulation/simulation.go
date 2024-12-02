package simulation

import (
	"fmt"
	"log"
	"time"

	agt "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/simulation/agent"
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/simulation/environment"
)

type Simulation struct {
	env     *envpkg.Environment
	agts    []envpkg.Agent
	objs    []envpkg.Object
	running bool
}

type SimulationJson struct {
	Agents  []interface{} `json:"agents"`
	Objects []interface{} `json:"objects"`
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
	simu.running = true
	log.Printf("Simulation started")

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
				simu.env.Lock()
				c <- simu.running
				time.Sleep(1 * time.Millisecond) // attente avant de relancer l'agent
				<-c
				simu.env.Unlock()
			}
		}(agt)
	}
}

func (simu *Simulation) Stop() {
	simu.env.Lock()
	simu.running = false
	simu.env.Unlock()
}

// Intention d'en faire un microservice
func (simu *Simulation) log() {
	// TODO
}

// Intention d'en faire un microservice
func (simu *Simulation) print() {
	startTime := time.Now()
	for simu.running {
		fmt.Printf("\rRunning simulation for %vms...  - ", time.Since(startTime).Milliseconds())
		time.Sleep(time.Second / 60) // 60 fps
	}
	log.Printf("Simulation stopped")
}

func (simu Simulation) ToJsonObj() SimulationJson {
	agents := []interface{}{}
	objects := []interface{}{}

	for _, agt := range simu.agts {
		agents = append(agents, agt.ToJsonObj())
	}

	for _, obj := range simu.objs {
		objects = append(objects, obj.ToJsonObj())
	}

	return SimulationJson{Agents: agents, Objects: objects}
}
