package simulation

import (
	"fmt"
	"log"
	"time"

	agt "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/agent"
	envpkg "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/environment"
	obj "gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation/object"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

type Simulation struct {
	env     *envpkg.Environment
	agts    []envpkg.IAgent
	objs    []envpkg.IObject
	running bool
}

type SimulationJson struct {
	Agents      []interface{} `json:"agents"`
	Objects     []interface{} `json:"objects"`
	Environment interface{}   `json:"environment"`
}

func NewSimulation(nagt int, nobj int) *Simulation {
	simu := &Simulation{}
	env := envpkg.NewEnvironment([]envpkg.IAgent{}, []envpkg.IObject{})

	simu.env = env

	mapDimension := utils.GetMapDimension()

	for i := 0; i < nagt; i++ {
		// création de l'agent
		id := fmt.Sprintf("ExAgent #%d", i)
		syncChan := make(chan bool)
		pos := envpkg.NewPosition(i, i, mapDimension, mapDimension)
		agt := agt.NewExAgent(id, pos, simu.env, syncChan)

		// ajout de l'agent à la simulation
		simu.agts = append(simu.agts, agt)

		// ajout de l'agent à l'environnement
		simu.env.AddAgent(agt)
	}

	for i := 0; i < nobj; i++ {
		// création de l'objet
		id := fmt.Sprintf("Flower #%d", i)
		pos := envpkg.NewPosition(nobj-i, i, mapDimension, mapDimension)
		obj := obj.NewFlower(id, pos)

		// ajout de l'objet à la simulation
		simu.objs = append(simu.objs, obj)

		// ajout de l'objet à l'environnement
		simu.env.AddObject(obj)
	}

	// Création d'une ruche
	nhive := 1
	for i := 0; i < nhive; i++ {
		// création de l'objet
		id := fmt.Sprintf("Hive #%d", i)
		pos := envpkg.NewPosition(9, 8, mapDimension, mapDimension)
		obj := obj.NewHive(id, pos, 0, 0, 0, 10)

		// ajout de l'objet à la simulation
		simu.objs = append(simu.objs, obj)

		// ajout de l'objet à l'environnement
		simu.env.AddObject(obj)
	}

	return simu
}

func (simu Simulation) IsRunning() bool {
	return simu.running
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
		go func(agt envpkg.IAgent) {
			for {
				c := agt.GetSyncChan()
				simu.env.Lock()
				c <- simu.running
				time.Sleep(100 * time.Millisecond) // attente avant de relancer l'agent
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

	return SimulationJson{Agents: agents, Objects: objects, Environment: simu.env.ToJsonObj()}
}
