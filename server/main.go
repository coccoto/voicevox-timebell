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
}

func listen() {
	fmt.Println(fmt.Sprintf("{comment}"))

	if err := http.ListenAndServe(SERVER_PORT, nil); err != nil {
		fmt.Println(fmt.Sprintf("{comment}"))
		os.Exit(1)
	}
}
