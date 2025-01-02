package simulation

import (
	"fmt"
	"log"
	"math/rand"
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

func NewSimulation(nbees int, nflowers int, nhornets int, maWs *websocket.Conn) *Simulation {
	simu := &Simulation{}
	simu.ws = maWs
	env := envpkg.NewEnvironment([]envpkg.IAgent{}, []envpkg.IObject{})
	simu.env = env
	//On récupère la webSocket
	mapDimension := utils.GetMapDimension()

	// Création d'une ruche
	nhive := 1
	hiveList := make([]*obj.Hive, nhive)
	for i := 0; i < nhive; i++ {
		// création de l'objet
		id := fmt.Sprintf("Hive #%d", i)
		x, y := rand.Intn(mapDimension), rand.Intn(mapDimension)
		pos := envpkg.NewPosition(x, y, mapDimension, mapDimension)
		hive := obj.NewHive(id, pos, 0, 0, 0, 10)
		// ajout de l'objet à la simulation
		simu.objs = append(simu.objs, hive)
		// ajout de l'objet à l'environnement
		simu.env.AddObject(hive)
		hiveList[i] = hive
	}

	nbPatches := utils.GetNumberFlowerPatches()
	nbFlowersPerPatch := nflowers / nbPatches
	nflowersLeft := nflowers
	for i := 0; i < nbPatches; i++ {
		centerOfPatchX := rand.Intn(mapDimension)
		centerOfPatchY := rand.Intn(mapDimension)
		// We create a zone of 2/3rd of the number of flowers
		offset := int(float64(nbFlowersPerPatch) * 2.0 / 3.0)
		// To avoid having flowers outside the map and delay due to unlucky guesses from machine
		// we generate a list of available positions
		availablePositions := make([]*envpkg.Position, 0)
		for x := centerOfPatchX - offset; x < centerOfPatchX+offset; x++ {
			for y := centerOfPatchY - offset; y < centerOfPatchY+offset; y++ {
				if simu.env.IsValidPosition(x, y) && simu.env.GetAt(x, y) == nil {
					availablePositions = append(availablePositions, envpkg.NewPosition(x, y, mapDimension, mapDimension))
				}
			}
		}
		// Now we can generate the flowers
		for j := 0; j < nbFlowersPerPatch || (i == nbPatches-1 && nflowersLeft > 0); j++ {
			// création de l'objet
			id := fmt.Sprintf("Flower #%d", i*nbFlowersPerPatch+j)
			// We create a zone of 2/3rd of the number of flowers
			selectedPosIdx := rand.Intn(len(availablePositions))
			pos := availablePositions[selectedPosIdx]
			availablePositions = append(availablePositions[:selectedPosIdx], availablePositions[selectedPosIdx+1:]...)
			flower := obj.NewFlower(id, pos)
			// ajout de l'objet à la simulation
			simu.objs = append(simu.objs, flower)
			// ajout de l'objet à l'environnement
			simu.env.AddObject(flower)
			nflowersLeft--
		}
	}

	// for i := 0; i < nbees; i++ {
	// 	// création de l'agent
	// 	id := fmt.Sprintf("ExAgent #%d", i)
	// 	syncChan := make(chan bool)
	// 	pos := envpkg.NewPosition(i, i, mapDimension, mapDimension)
	// 	agt := agt.NewExAgent(id, pos, simu.env, syncChan)

	// 	// ajout de l'agent à la simulation
	// 	simu.agts = append(simu.agts, agt)

	// 	// ajout de l'agent à l'environnement
	// 	simu.env.AddAgent(agt)
	// }

	maxNectar := utils.GetMaxNectar()
	for i := 0; i < nbees; i++ {
		// Creating a bee
		id := fmt.Sprintf("Bee #%d", i)
		syncChan := make(chan bool)
		hive := hiveList[rand.Intn(nhive)]
		agt := agt.NewBeeAgent(id, simu.env, syncChan, rand.Intn(2)+1, hive, time.Now(), maxNectar, agt.Worker)
		// ajout de l'agent à la simulation
		simu.agts = append(simu.agts, agt)
		// ajout de l'agent à l'environnement
		simu.env.AddAgent(agt)
	}

	for i := 0; i < nhornets; i++ {
		// Creating a hornet
		id := fmt.Sprintf("Hornet #%d", i)
		syncChan := make(chan bool)
		agt := agt.NewHornetAgent(id, simu.env, syncChan, rand.Intn(2)+1)
		// ajout de l'agent à la simulation
		simu.agts = append(simu.agts, agt)
		// ajout de l'agent à l'environnement
		//simu.env.AddAgent(agt) // Pas encore car pas spawn
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
	for simu.IsRunning() {
		for i, agt := range simu.agts {
			c := agt.GetSyncChan()
			simu.env.Lock()
			c <- true
			isAlive := <-c
			// If dead, remove agent from simulation
			if !isAlive {
				fmt.Printf("{{SIMULATION}} - [%s] is dead\n", agt.ID())
				simu.agts = append(simu.agts[:i], simu.agts[i+1:]...)
			}
			simu.env.Unlock()
			time.Sleep(time.Second / 100) // 100 Tour / Sec
		}
	}

	// Arrêt des agents
	for _, agt := range simu.agts {
		c := agt.GetSyncChan()
		c <- false
		<-c
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
		time.Sleep(time.Second / 30) // 60 fps
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

	simu.env.Lock()
	defer simu.env.Unlock()

	return SimulationJson{Agents: agents, Objects: objects, Environment: simu.env.ToJsonObj()}
}
