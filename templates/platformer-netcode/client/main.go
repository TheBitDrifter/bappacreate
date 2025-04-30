package main

import (
	"log"

	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_clientsystems"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_rendersystems"

	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/actions"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/scenes"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/sharedclient"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/sharedclient/assets"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/sharedclient/clientsystems"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/sharedclient/rendersystems"

	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/coresystems"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	log.Println("Starting Networked Client...")

	client := coldbrew.NewNetworkClient(
		sharedclient.RESOLUTION_X,
		sharedclient.RESOLUTION_Y,
		sharedclient.MAX_SPRITES_CACHED,
		sharedclient.MAX_SOUNDS_CACHED,
		sharedclient.MAX_SCENES_CACHED,
		assets.FS,
	)

	client.SetDeserCallback(Derser)

	client.SetLocalAssetPath("../sharedclient/assets/")

	// Client Settings
	client.SetTitle("Platformer LDTK Template (Networked)")
	client.SetResizable(true)
	client.SetMinimumLoadTime(30)

	log.Println("Registering Scene One...")
	err := client.RegisterScene(
		scenes.SceneOne.Name,
		scenes.SceneOne.Width,
		scenes.SceneOne.Height,
		scenes.SceneOne.Plan,
		rendersystems.DefaultRenderSystems,
		clientsystems.DefaultClientSystemsNetworked,
		// []blueprint.CoreSystem{}, <- for non interpolated use this instead
		coresystems.DefaultCoreSystems, // <- for interpolation use this
		scenes.SceneOne.Preload...,
	)
	if err != nil {
		log.Fatalf("Failed to register Scene One: %v", err)
	}

	// Register Global Systems
	log.Println("Registering Global Systems...")
	client.RegisterGlobalRenderSystem(
		coldbrew_rendersystems.GlobalRenderer{},
		&coldbrew_rendersystems.DebugRenderer{},
	)
	client.RegisterGlobalClientSystem(
		&coldbrew_clientsystems.InputSenderSystem{},
		coldbrew_clientsystems.InputBufferSystem{},
		&coldbrew_clientsystems.CameraSceneAssignerSystem{},
	)

	log.Println("Activating Camera...")
	_, err = client.ActivateCamera()
	if err != nil {
		log.Fatalf("Failed to activate camera: %v", err)
	}

	log.Println("Activating Input Receiver and Mapping Keys...")
	receiver1, err := client.ActivateReceiver()
	if err != nil {
		log.Fatalf("Failed to activate receiver: %v", err)
	}
	receiver1.RegisterKey(ebiten.KeySpace, actions.Jump)
	receiver1.RegisterKey(ebiten.KeyW, actions.Jump)
	receiver1.RegisterKey(ebiten.KeyA, actions.Left)
	receiver1.RegisterKey(ebiten.KeyD, actions.Right)
	receiver1.RegisterKey(ebiten.KeyS, actions.Down)

	log.Printf("Connecting to Drip server at %s...", sharedclient.SERVER_ADDRESS)
	err = client.Connect(sharedclient.SERVER_ADDRESS)
	if err != nil {
		log.Fatalf("Failed to connect to server '%s': %v", sharedclient.SERVER_ADDRESS, err)
	}
	defer func() {
		log.Println("Disconnecting from server...")
		client.Disconnect()
	}()
	log.Println("Connected successfully.")

	log.Println("Starting Ebiten game loop (blocking)...")
	if err := client.Start(); err != nil {
		log.Fatalf("Client exited with error: %v", err)
	}

	log.Println("Client shutdown complete.")
}
