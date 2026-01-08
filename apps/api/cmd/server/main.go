package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/config"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/httpapi"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create server
	server, err := httpapi.NewServer(cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	defer server.Close()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("shutting down...")
		server.Close()
		os.Exit(0)
	}()

	// Start server
	log.Printf("api listening on %s (mode: %s)", cfg.HTTPAddr, cfg.AppMode)
	if err := server.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
