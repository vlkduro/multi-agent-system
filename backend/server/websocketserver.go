package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/simulation"
)

type WebSocketServer struct {
	port       int
	simulation *simulation.Simulation
	running    bool
}

func MakeWebSocketServer(port int) *WebSocketServer {
	return &WebSocketServer{port: port, simulation: nil, running: false}
}

func (server *WebSocketServer) newSimulation(w http.ResponseWriter, _ *http.Request) {
	if server.simulation != nil {
		server.simulation.Stop()
	}
	server.simulation = simulation.NewSimulation(10, 10)
	w.WriteHeader(http.StatusOK)
}

func (server *WebSocketServer) launchSimulation(w http.ResponseWriter, _ *http.Request) {
	if server.simulation == nil {
		server.simulation = simulation.NewSimulation(10, 10)
	}

	if server.simulation.IsRunning() {
		w.WriteHeader(http.StatusConflict)
		return
	}

	go server.simulation.Run()
	server.running = true
	w.WriteHeader(http.StatusOK)
}

func (server *WebSocketServer) stopSimulation(w http.ResponseWriter, _ *http.Request) {
	if server.simulation == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	server.simulation.Stop()
	w.WriteHeader(http.StatusOK)
}

func (server *WebSocketServer) connectToSimulation(w http.ResponseWriter, r *http.Request) {
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

	for server.running {
		time.Sleep(time.Second / 10) // 60 tps
		if server.simulation == nil {
			continue
		}
		if err := conn.WriteJSON(server.simulation.ToJsonObj()); err != nil {
			log.Println(err)
			break
		}
	}
	log.Println("WebSocket server stopped")
}

func (server *WebSocketServer) setupRoutes() int {
	http.HandleFunc("/new", server.newSimulation)
	http.HandleFunc("/start", server.launchSimulation)
	http.HandleFunc("/stop", server.stopSimulation)

	http.HandleFunc("/connect", server.connectToSimulation)

	return http.StatusOK
}

func (server *WebSocketServer) LaunchServer() {
	server.setupRoutes()
	server.running = true

	log.Printf("WebSocket server starting on port %d", server.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (server *WebSocketServer) StopServer() {
	log.Printf("WebSocket server stopping")
	server.running = false
}
