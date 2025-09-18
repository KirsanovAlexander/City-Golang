package main

import (
	"log"
	"os"

	"city/internal/router"
	"city/internal/storage"
)

func main() {
	store := storage.NewMemoryStore()
	r := router.SetupRouter(store)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
