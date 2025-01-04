package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/simulation"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/utils"
)

type WebSocketServer struct {
	port       int
	simulation *simulation.Simulation
	running    bool
}

func MakeWebSocketServer(port int) *WebSocketServer {
	return &WebSocketServer{port: port, simulation: nil, running: false}
}

func (server *WebSocketServer) newSimulation(conn *websocket.Conn) {
	if server.simulation != nil {
		fmt.Print("old running")
		server.simulation.Stop()
	}
	nAgts := utils.GetNumberBees()
	nObjs := utils.GetNumberObjects()
	nhornets := utils.GetNumberHornets()
	server.simulation = simulation.NewSimulation(nAgts, nObjs, nhornets, conn)
}

func (server *WebSocketServer) launchSimulation(
	conn *websocket.Conn,
) {
	if server.simulation == nil {
		nAgts := utils.GetNumberBees()
		nObjs := utils.GetNumberObjects()
		nhornets := utils.GetNumberHornets()
		server.simulation = simulation.NewSimulation(nAgts, nObjs, nhornets, conn)
	}

	if server.simulation.IsRunning() {
		fmt.Print("Allready running")
		return
	}

	go server.simulation.Run(conn)
}

func (server *WebSocketServer) stopSimulation() {
	if server.simulation == nil {
		return
	}
	server.simulation.Stop()
}

func (server *WebSocketServer) startSimulation(w http.ResponseWriter, r *http.Request) {
	// Setting up the websocket with default parameters
	upgrader := websocket.Upgrader{}
	// Accept all origins
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	for {
		var message string = ""
		if conn == nil {
			conn, err = upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			// Lire un message du client
			_, text, err := conn.ReadMessage()
			message = string(text)
			if err != nil {
				log.Println("Conn lost error:", err)
				server.stopSimulation()
				conn.Close()
				conn = nil
			}
		}
		switch message {
		case "start":
			server.launchSimulation(conn)
		case "stop":
			server.stopSimulation()
		case "new":
			if err != nil {
				log.Println(err)
				return
			}
			server.newSimulation(conn)
		case "bye":
			conn.Close()
			conn = nil
		}
	}
}

func (server *WebSocketServer) LaunchServer() {
	http.HandleFunc("/ws/", server.startSimulation)

	log.Printf("WebSocket server starting on port %d", server.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (server *WebSocketServer) StopServer() {
	log.Printf("WebSocket server stopping")
}
