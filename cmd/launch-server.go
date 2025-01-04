package main

import (
	"fmt"

	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/backend/server"
	"gitlab.utc.fr/bidauxal/ai30_valakou_martins_chartier_bidaux/frontend"
)

func main() {
	srv := server.MakeWebSocketServer(8080)
	go srv.LaunchServer()
	go frontend.LaunchFileServiceServer()
	fmt.Scanln()
	srv.StopServer()
}
