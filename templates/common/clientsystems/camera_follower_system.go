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
	playersWithCamera := warehouse.Factory.NewQuery()
	playersWithCamera.And(
		spatial.Components.Position,
		input.Components.ActionBuffer,
		client.Components.CameraIndex,
	)

	playerCursor := scene.NewCursor(playersWithCamera)

	for range playerCursor.Next() {
		playerPos := spatial.Components.Position.GetFromCursor(playerCursor)

		camIndex := int(*client.Components.CameraIndex.GetFromCursor(playerCursor))
		cam := cli.Cameras()[camIndex]

		_, cameraScenePosition := cam.Positions()
		centerX := float64(cam.Surface().Bounds().Dx()) / 2
		centerY := float64(cam.Surface().Bounds().Dy()) / 2

		centeredPlayerX := playerPos.X
		centeredPlayerY := playerPos.Y
		centeredCameraX := cameraScenePosition.X + centerX
		centeredCameraY := cameraScenePosition.Y + centerY

		diffX := centeredPlayerX - centeredCameraX
		diffY := centeredPlayerY - centeredCameraY

		deadzoneX := 60.0 // horizontal deadzone in pixels
		deadzoneY := 60.0 // vertical deadzone in pixels

		targetX := cameraScenePosition.X
		targetY := cameraScenePosition.Y

		// Only move camera if player outside deadzone
		if math.Abs(diffX) > deadzoneX {
			// Adjust target position to keep player at edge of deadzone X
			if diffX > 0 {
				targetX = centeredPlayerX - centerX - deadzoneX
			} else {
				targetX = centeredPlayerX - centerX + deadzoneX
			}
		}

		if math.Abs(diffY) > deadzoneY {
			// Adjust target position to keep player at edge of deadzone Y
			if diffY > 0 {
				targetY = centeredPlayerY - centerY - deadzoneY
			} else {
				targetY = centeredPlayerY - centerY + deadzoneY
			}
		}

		// Apply smooth lerping to camera movement
		cameraScenePosition.X = lerp(cameraScenePosition.X, targetX, 0.03)
		cameraScenePosition.Y = lerp(cameraScenePosition.Y, targetY, 0.03)

		lockCameraToSceneBoundaries(cam, scene, cameraScenePosition)
	}
	return nil
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// lockCameraToSceneBoundaries constrains camera position within scene boundaries
func lockCameraToSceneBoundaries(cam coldbrew.Camera, scene coldbrew.Scene, cameraPos *vector.Two) {
	sceneWidth := scene.Width()
	sceneHeight := scene.Height()
	camWidth, camHeight := cam.Dimensions()
	maxX := sceneWidth - camWidth
	maxY := sceneHeight - camHeight

	// Constrain X
	if cameraPos.X > float64(maxX) {
		cameraPos.X = float64(maxX)
	}
	if cameraPos.X < 0 {
		cameraPos.X = 0
	}
	// Constrain Y
	if cameraPos.Y > float64(maxY) {
		cameraPos.Y = float64(maxY)
	}
	if cameraPos.Y < 0 {
		cameraPos.Y = 0
	}
}
