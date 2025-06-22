package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	SERVER_PORT = ":8080"
)

func main() {
	setupRoutes()
	listen()
}

func setupRoutes() {
	http.HandleFunc("/api/alert", corsMiddleware(alertHandler))
	http.HandleFunc("/api/config", corsMiddleware(configHandler))
	http.HandleFunc("/api/speakers", corsMiddleware(voicevoxSpeakersHandler))
}

func listen() {
	fmt.Println("Server listening on port", SERVER_PORT)

	if err := http.ListenAndServe(SERVER_PORT, nil); err != nil {
		fmt.Println("Failed to listen:", err)
		os.Exit(1)
	}
}
