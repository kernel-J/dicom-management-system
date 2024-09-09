package main

import (
	"log"

	"ImageManagementService/config"
	"ImageManagementService/server"
)

func main() {
	// Load configuration
	cfg := config.New()
	cfg.Validate()

	// Start server
	srv := server.NewServer(cfg)
	log.Printf("Starting server on port %s\n", cfg.Port)
	srv.Run()
}