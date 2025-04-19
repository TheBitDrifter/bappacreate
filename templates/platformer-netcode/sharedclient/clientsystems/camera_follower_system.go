package clientsystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
)

type CameraFollowerSystem struct{}

func (CameraFollowerSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	netCli, isNet := cli.(coldbrew.NetworkClient)

	if isNet {
		id, ok := netCli.AssociatedEntityID()
		if !ok {
			return nil
		}

		pEn, err := scene.Storage().Entity(id)
		if err != nil {
			return nil
		}

		follow(
			netCli,
			scene,
			*spatial.Components.Position.GetFromEntity(pEn),
			int(*client.Components.CameraIndex.GetFromEntity(pEn)),
		)
		return nil
	}

	playersWithCamera := warehouse.Factory.NewQuery().And(
		spatial.Components.Position,
		input.Components.ActionBuffer,
		client.Components.CameraIndex,
	)
	playerCursor := scene.NewCursor(playersWithCamera)

	for range playerCursor.Next() {
		follow(
			cli,
			scene,
			*spatial.Components.Position.GetFromCursor(playerCursor),
			int(*client.Components.CameraIndex.GetFromCursor(playerCursor)),
		)
	}
	return nil
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func follow(cli coldbrew.LocalClient, scene coldbrew.Scene, playerPos spatial.Position, camIndex int) {
	cam := cli.Cameras()[camIndex]

	// Get the cameras local scene position
	_, cameraScenePosition := cam.Positions()
	centerX := float64(cam.Surface().Bounds().Dx()) / 2
	centerY := float64(cam.Surface().Bounds().Dy()) / 2

	// calculate the centered positions for BOTH player and camera
	centeredPlayerX := playerPos.X
	centeredPlayerY := playerPos.Y
	centeredCameraX := cameraScenePosition.X + centerX
	centeredCameraY := cameraScenePosition.Y + centerY

	// Calculate distance between camera center and player position
	diffX := centeredPlayerX - centeredCameraX
	diffY := centeredPlayerY - centeredCameraY

	// deadzone - camera only moves when player is outside of deadzone
	deadzoneX := 60.0 // horizontal deadzone in pixels
	deadzoneY := 60.0 // vertical deadzone in pixels

	// Target position starts at current camera position
	targetX := cameraScenePosition.X
	targetY := cameraScenePosition.Y

	// Only move camera if player outside deadzone
	if math.Abs(diffX) > deadzoneX {
		// Adjust target position to keep player at edge of deadzone
		if diffX > 0 {
			targetX = centeredPlayerX - centerX - deadzoneX
		} else {
			targetX = centeredPlayerX - centerX + deadzoneX
		}
	}

	if math.Abs(diffY) > deadzoneY {
		// Adjust target position to keep player at edge of deadzone
		if diffY > 0 {
			targetY = centeredPlayerY - centerY - deadzoneY
		} else {
			targetY = centeredPlayerY - centerY + deadzoneY
		}
	}

	// Apply smooth lerping to camera movement
	cameraScenePosition.X = lerp(cameraScenePosition.X, targetX, 0.03)
	cameraScenePosition.Y = lerp(cameraScenePosition.Y, targetY, 0.03)

	// Lock the camera to the scene boundaries
	lockCameraToSceneBoundaries(cam, scene, cameraScenePosition)
}

func lockCameraToSceneBoundaries(cam coldbrew.Camera, scene coldbrew.Scene, cameraPos *vector.Two) {
	sceneWidth := scene.Width()
	sceneHeight := scene.Height()
	camWidth, camHeight := cam.Dimensions()
	// Calculate maximum positions to keep camera within scene bounds
	maxX := sceneWidth - camWidth
	maxY := sceneHeight - camHeight
	// Constrain camera X position
	if cameraPos.X > float64(maxX) {
		cameraPos.X = float64(maxX)
	}
	if cameraPos.X < 0 {
		cameraPos.X = 0
	}
	// Constrain camera Y position
	if cameraPos.Y > float64(maxY) {
		cameraPos.Y = float64(maxY)
	}
	if cameraPos.Y < 0 {
		cameraPos.Y = 0
	}
}
