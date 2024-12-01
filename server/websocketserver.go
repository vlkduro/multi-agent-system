package server

import (
	"encoding/json"
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
}

func MakeWebSocketServer(port int) *WebSocketServer {
	return &WebSocketServer{port: port, simulation: nil}
}

func (server *WebSocketServer) launchSimulation(w http.ResponseWriter, _ *http.Request) {
	if server.simulation != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	server.simulation = simulation.NewSimulation(10, 10)
	go server.simulation.Run()
	w.WriteHeader(http.StatusOK)
}

func (server *WebSocketServer) connectToSimulation(w http.ResponseWriter, r *http.Request) {
	// Setting up the websocket with default parameters
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		w.WriteHeader(http.StatusOK)
		serial, err := json.Marshal(server.simulation.ToJsonObj())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err.Error())
			return
		}
		if err := conn.WriteMessage(messageType, serial); err != nil {
			log.Println(err)
			break
		}
		time.Sleep(time.Second / 60) // 60 tps
	}
}

func (server *WebSocketServer) setupRoutes() int {
	http.HandleFunc("/start", server.launchSimulation)

	http.HandleFunc("/connect", server.connectToSimulation)

	return http.StatusOK
}

func (server *WebSocketServer) LaunchServer() {
	server.setupRoutes()

	log.Printf("WebSocket server starting on port %d", server.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
