package frontend

import (
	"fmt"
	"log"
	"net/http"
)

func LaunchFileServiceServer() {
	// Définir les options de ligne de commande
	port := 8000
	dir := "./frontend"

	// Créer un gestionnaire de fichiers statiques
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", http.StripPrefix("/", fs))

	// Lancer le serveur HTTP
	log.Printf("WebSocket server starting on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
