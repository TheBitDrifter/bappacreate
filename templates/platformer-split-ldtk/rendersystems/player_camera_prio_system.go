package rendersystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_rendersystems"
)

// When split screen side scrolling its better to render the player at the highest level for THIER RESPECTIVE CAMERA
type PlayerCameraPriorityRenderer struct{}

func (PlayerCameraPriorityRenderer) Render(scene coldbrew.Scene, screen coldbrew.Screen, cli coldbrew.LocalClient) {
	for _, cam := range cli.ActiveCamerasFor(scene) {
		if !cli.Ready(cam) {
			continue
		}
		playerCursor := scene.NewCursor(blueprint.Queries.ActionBuffer)

		for range playerCursor.Next() {

			camIndex := client.Components.CameraIndex.GetFromCursor(playerCursor)
			// Skip entities that match the current camera's index
			if int(*camIndex) == cam.Index() {
				continue
			}
			coldbrew_rendersystems.RenderEntityFromCursor(playerCursor, cam, scene.CurrentTick())
		}

		for range playerCursor.Next() {
			camIndex := client.Components.CameraIndex.GetFromCursor(playerCursor)

			// Only render entities that match the current camera's index
			if int(*camIndex) == cam.Index() {
				coldbrew_rendersystems.RenderEntityFromCursor(playerCursor, cam, scene.CurrentTick())
			}
		}
		cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
	}
}
