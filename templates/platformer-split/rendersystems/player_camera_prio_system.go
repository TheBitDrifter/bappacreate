package rendersystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_rendersystems"
)

type PlayerCameraPriorityRenderer struct{}

// When split screen side scrolling its better to render the player at the highest level for THIER RESPECTIVE CAMERA
func (PlayerCameraPriorityRenderer) Render(scene coldbrew.Scene, screen coldbrew.Screen, cameraUtility coldbrew.CameraUtility) {
	// Loop through active cameras for the scene
	for _, cam := range cameraUtility.ActiveCamerasFor(scene) {
		// If it ain't ready chill out!
		if !cameraUtility.Ready(cam) {
			continue
		}
		// First, render all non-matching entities (the players unrelated to the current camera)
		playerCursor := scene.NewCursor(blueprint.Queries.InputBuffer)
		for range playerCursor.Next() {
			camIndex := client.Components.CameraIndex.GetFromCursor(playerCursor)
			// Skip entities that match the current camera's index
			if int(*camIndex) == cam.Index() {
				continue
			}
			coldbrew_rendersystems.RenderEntityFromCursor(playerCursor, cam, scene.CurrentTick())
		}

		// Then render matching entities last (the player that 'owns' the camera)
		playerCursor = scene.NewCursor(blueprint.Queries.InputBuffer)
		for range playerCursor.Next() {
			camIndex := client.Components.CameraIndex.GetFromCursor(playerCursor)
			// Only render entities that match the current camera's index
			if int(*camIndex) == cam.Index() {
				coldbrew_rendersystems.RenderEntityFromCursor(playerCursor, cam, scene.CurrentTick())
			}
		}
		cam.PresentToScreen(screen, 3)
	}
}
