package simulation

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
	ws      *websocket.Conn
	sync.Mutex
}

type SimulationJson struct {
	Agents      []interface{} `json:"agents"`
	Objects     []interface{} `json:"objects"`
	Environment interface{}   `json:"environment"`
}

func NewSimulation(nagt int, nobj int, maWs *websocket.Conn) *Simulation {
	simu := &Simulation{}
	simu.ws = maWs
	env := envpkg.NewEnvironment([]envpkg.IAgent{}, []envpkg.IObject{})
	simu.env = env
	//On récupère la webSocket
	mapDimension := utils.GetMapDimension()

	for i := 0; i < nagt; i++ {
		// création de l'agent
		id := fmt.Sprintf("Agent #%d", i)
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
		id := fmt.Sprintf("Object #%d", i)
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
	simu.sendState()
	return simu
}

func (simu *Simulation) IsRunning() bool {
	simu.Lock()
	defer simu.Unlock()
	return simu.running
}

func (simu *Simulation) Run(maWs *websocket.Conn) {
	simu.Lock()
	simu.running = true
	if maWs != simu.ws {
		simu.ws = maWs
	}
	log.Printf("Simulation started")
	simu.Unlock()

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
			c := agt.GetSyncChan()
			for {
				simu.env.Lock()
				c <- simu.IsRunning()
				time.Sleep(1 * time.Millisecond) // attente avant de relancer l'agent
				r := <-c
				simu.env.Unlock()
				if !r {
					break
				}
			}
		}(agt)
	}
}

func (simu *Simulation) Stop() {
	simu.Lock()
	defer simu.Unlock()
	simu.running = false
	simu.ws = nil
	log.Printf("\nSimulation stopped\n")
}

// Intention d'en faire un microservice
func (simu *Simulation) log() {
	if (simu.ws) == nil {
		return
	}
	for simu.IsRunning() {
		simu.sendState()
		time.Sleep(time.Second / 60) // 60 fps
	}
}

func (simu *Simulation) sendState() {
	err := (*simu.ws).WriteJSON(simu.ToJsonObj())
	if err != nil {
		log.Fatal("Log micro service: ", err)
	}
}

// Intention d'en faire un microservice
func (simu *Simulation) print() {
	startTime := time.Now()
	for simu.IsRunning() {
		fmt.Printf("\rRunning simulation for %vms...  - ", time.Since(startTime).Milliseconds())
		time.Sleep(time.Second / 60) // 60 fps
	}
}

func (simu *Simulation) ToJsonObj() SimulationJson {
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
