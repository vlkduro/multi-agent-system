package main

import (
	"fmt"

	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/server"
)

func main() {
	server := server.MakeWebSocketServer(8080)
	go server.LaunchServer()
	fmt.Scanln()
	server.StopServer()
}
