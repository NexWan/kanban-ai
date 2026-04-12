package main

import (
	"log"
	"net/http"

	"agent-kanban-api/internal/app"
)

func main() {
	addr := ":8080"
	router := app.NewRouter()

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
