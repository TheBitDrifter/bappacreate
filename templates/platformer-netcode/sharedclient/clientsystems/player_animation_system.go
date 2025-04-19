package clientsystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/animations"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/components"
)

type PlayerAnimationSystem struct{}

func (PlayerAnimationSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	cursor := scene.NewCursor(blueprint.Queries.ActionBuffer)

	for range cursor.Next() {
		bundle := client.Components.SpriteBundle.GetFromCursor(cursor)
		spriteBlueprint := &bundle.Blueprints[0] // Player entities currently only use one (first) blueprint

		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		grounded, onGround := components.OnGroundComponent.GetFromCursorSafe(cursor)
		if grounded {
			grounded = scene.CurrentTick()-onGround.LastTouch <= 3 // Slight tolerance for netcode
		}

		// Player is moving horizontal and grounded (running)
		if math.Abs(dyn.Vel.X) > 20 && grounded {
			spriteBlueprint.TryAnimation(animations.RunAnimation)

			// Player is moving down and not grounded (falling)
		} else if dyn.Vel.Y > 0 && !grounded {
			spriteBlueprint.TryAnimation(animations.FallAnimation)

			// Player is moving up and not grounded (jumping)
		} else if dyn.Vel.Y <= 0 && !grounded {
			spriteBlueprint.TryAnimation(animations.JumpAnimation)

			// Default: player is idle
		} else {
			spriteBlueprint.TryAnimation(animations.IdleAnimation)
		}
	}
	return nil
}
