package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappacreate/templates/topdown/actions"
	"github.com/TheBitDrifter/bappacreate/templates/topdown/components"
)

type PlayerMovementSystem struct{}

func (sys PlayerMovementSystem) Run(scene blueprint.Scene, dt float64) error {
	cursor := scene.NewCursor(blueprint.Queries.ActionBuffer)

	for range cursor.Next() {
		incomingInputs := input.Components.ActionBuffer.GetFromCursor(cursor)
		pos := spatial.Components.Position.GetFromCursor(cursor)
		direction := spatial.Components.Direction.GetFromCursor(cursor)
		direction8 := components.DirectionEightComponent.GetFromCursor(cursor)

		// Process horizontal movement
		_, pressedLeft := incomingInputs.ConsumeAction(actions.Left)
		_, pressedRight := incomingInputs.ConsumeAction(actions.Right)

		// Process vertical movement
		_, pressedUp := incomingInputs.ConsumeAction(actions.Up)
		_, pressedDown := incomingInputs.ConsumeAction(actions.Down)

		// Calculate movement vector
		moveX := 0.0
		moveY := 0.0

		if pressedLeft {
			moveX -= 1
		}
		if pressedRight {
			moveX += 1
		}
		if pressedUp {
			moveY -= 1
		}
		if pressedDown {
			moveY += 1
		}

		// Normalize diagonal movement
		if moveX != 0 && moveY != 0 {
			// Calculate vector length
			length := math.Sqrt(moveX*moveX + moveY*moveY)

			// Normalize the vector
			moveX = moveX / length
			moveY = moveY / length
		}

		// Apply movement
		pos.X += moveX
		pos.Y += moveY

		// Set direction components based on movement
		movingHorizontal := moveX != 0
		movingVertical := moveY != 0

		// Update cardinal direction component
		if moveX < 0 {
			direction.SetLeft()
		} else if moveX > 0 {
			direction.SetRight()
		}
		// Set the eight-direction component based on combined input
		if movingVertical || movingHorizontal {
			if pressedUp && pressedRight {
				direction8.SetRightUp()
			} else if pressedUp && pressedLeft {
				direction8.SetLeftUp()
			} else if pressedDown && pressedRight {
				direction8.SetRightDown()
			} else if pressedDown && pressedLeft {
				direction8.SetLeftDown()
			} else if pressedUp {
				direction8.SetUp()
			} else if pressedDown {
				direction8.SetDown()
			} else if pressedLeft {
				direction8.SetLeft()
			} else if pressedRight {
				direction8.SetRight()
			}
		}

		playerEn, err := cursor.CurrentEntity()
		if err != nil {
			return err
		}

		// Track movement
		if movingVertical || movingHorizontal {
			playerEn.EnqueueAddComponent(components.IsMovingComponent)
		} else {
			playerEn.EnqueueRemoveComponent(components.IsMovingComponent)
		}
	}
	return nil
}
