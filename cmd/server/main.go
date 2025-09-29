package main

import (
	"log"
	"net/http"
	"os"

	"city/internal/router"
	"city/internal/storage"
)

// @title City Management API
// @version 1.0
// @description API for managing a city simulation game
// @host localhost:8080
// @BasePath /
func main() {
	store := storage.NewMemoryStore()
	r := router.SetupRouter(store)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
