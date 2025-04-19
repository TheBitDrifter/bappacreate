package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/TheBitDrifter/bappa/drip"
	"github.com/TheBitDrifter/bappa/drip/drip_seversystems"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/coresystems"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/scenes"
)

func main() {
	drip.Callbacks.NewConnectionCreateEntity = NewConnectionEntityCreate
	drip.Callbacks.Serialize = SerializeCallback

	config := drip.DefaultServerConfig()

	server := drip.NewServer(config, drip_seversystems.ActionBufferSystem{})

	// Register a scene
	log.Println("Registering scene:", scenes.SceneOne.Name)
	err := server.RegisterScene(
		scenes.SceneOne.Name,
		scenes.SceneOne.Width,
		scenes.SceneOne.Height,
		scenes.SceneOne.Plan,
		coresystems.DefaultCoreSystems,
	)
	if err != nil {
		log.Fatalf("Failed to register scene: %v", err)
	}

	// Start the server
	log.Println("Starting server...")
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)

	// Notify the channel for specific signals (Interrupt, Terminate)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Server running. Press Ctrl+C to stop.")

	// Block execution until a signal is received on the 'quit' channel
	<-quit

	// Initiate shutdown
	log.Println("Shutting down server...")
	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	} else {
		log.Println("Server stopped gracefully.")
	}
}
