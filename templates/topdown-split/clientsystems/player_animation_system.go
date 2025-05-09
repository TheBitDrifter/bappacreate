package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/animations"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/components"
)

type PlayerAnimationSystem struct{}

func (PlayerAnimationSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	const PLAYER_IDLE_SHEET_INDEX = 0
	const PLAYER_WALK_SHEET_INDEX = 1

	// Iterate through players
	cursor := scene.NewCursor(blueprint.Queries.ActionBuffer)
	for range cursor.Next() {
		direction8 := components.DirectionEightComponent.GetFromCursor(cursor)
		bundle := client.Components.SpriteBundle.GetFromCursor(cursor)
		playerMoving := components.IsMovingComponent.CheckCursor(cursor)
		spriteBlueprint := &bundle.Blueprints[PLAYER_IDLE_SHEET_INDEX]

		// Update sheet based on movement
		if playerMoving {
			spriteBlueprint = &bundle.Blueprints[PLAYER_IDLE_SHEET_INDEX]
			spriteBlueprint.Deactivate()
			spriteBlueprint = &bundle.Blueprints[PLAYER_WALK_SHEET_INDEX]
			spriteBlueprint.Activate()
		} else {
			spriteBlueprint = &bundle.Blueprints[PLAYER_WALK_SHEET_INDEX]
			spriteBlueprint.Deactivate()
			spriteBlueprint = &bundle.Blueprints[PLAYER_IDLE_SHEET_INDEX]
			spriteBlueprint.Activate()
		}

		// Based on the DirectionEight and Direction we pick the sprite and flip accordingly
		// DirectionEight informs us of the animation (sheet) to use
		// Direction (left/right) informs the GlobalRenderer when to flip the sprite
		if direction8.IsDown() {
			spriteBlueprint.TryAnimation(animations.Down)
		} else if direction8.IsUp() {
			spriteBlueprint.TryAnimation(animations.Up)
		} else if direction8.IsRight() || direction8.IsLeft() {
			spriteBlueprint.TryAnimation(animations.Side)
		} else if direction8.IsRightDown() || direction8.IsLeftDown() {
			spriteBlueprint.TryAnimation(animations.DownSide)
		} else if direction8.IsRightUp() || direction8.IsLeftUp() {
			spriteBlueprint.TryAnimation(animations.UpSide)
		}
	}
	return nil
}
