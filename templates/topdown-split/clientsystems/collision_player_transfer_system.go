package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/components"
)

type CollisionPlayerTransferSystem struct{}

type playerTransfer struct {
	target       string
	playerEntity warehouse.Entity
}

// System handles transferring players when they collide with PlayerSceneTransfer entities
func (CollisionPlayerTransferSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	var pending []playerTransfer

	collisionTransferQuery := warehouse.Factory.NewQuery().And(
		spatial.Components.Shape,
		components.PlayerSceneTransferComponent,
	)

	playerWithShapeQuery := warehouse.Factory.NewQuery().And(
		spatial.Components.Shape,
		input.Components.ActionBuffer,
	)

	collisionTransferCursor := scene.NewCursor(collisionTransferQuery)
	playerWithShapeCursor := scene.NewCursor(playerWithShapeQuery)

	for range collisionTransferCursor.Next() {
		transferPos := spatial.Components.Position.GetFromCursor(collisionTransferCursor)
		transferCollider := spatial.Components.Shape.GetFromCursor(collisionTransferCursor)

		for range playerWithShapeCursor.Next() {

			playerPos := spatial.Components.Position.GetFromCursor(playerWithShapeCursor)
			playerCollider := spatial.Components.Shape.GetFromCursor(playerWithShapeCursor)

			if ok, _ := spatial.Detector.Check(*playerCollider, *transferCollider, playerPos, transferPos); ok {
				sceneTransfer := components.PlayerSceneTransferComponent.GetFromCursor(collisionTransferCursor)
				playerEn, err := playerWithShapeCursor.CurrentEntity()
				if err != nil {
					return err
				}

				transfer := playerTransfer{
					target:       sceneTransfer.Dest,
					playerEntity: playerEn,
				}
				// Enqueue transfer
				pending = append(pending, transfer)

				// Update the player pos
				playerPos := spatial.Components.Position.GetFromCursor(playerWithShapeCursor)
				playerPos.X = sceneTransfer.X
				playerPos.Y = sceneTransfer.Y

				// Update the camera pos
				camIndex := int(*client.Components.CameraIndex.GetFromCursor(playerWithShapeCursor))
				cam := cli.Cameras()[camIndex]

				_, cameraScenePosition := cam.Positions()
				centerX := float64(cam.Surface().Bounds().Dx()) / 2
				centerY := float64(cam.Surface().Bounds().Dy()) / 2

				cameraScenePosition.X = playerPos.X - centerX
				cameraScenePosition.Y = playerPos.Y - centerY

			}
		}

	}

	// Process transfers after loop
	for _, transfer := range pending {
		_, err := cli.ActivateSceneByName(transfer.target, transfer.playerEntity)
		if err != nil {
			return err
		}
	}

	return nil
}
