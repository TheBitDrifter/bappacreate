package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/actions"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/components"
)

const (
	SPEED_X            = 120.0 // Player's horizontal movement speed
	SNAP_FORCE         = 40.0  // Downward force to keep player attached to slopes
	JUMP_FORCE         = 320.0 // Upward force applied when jumping
	COYOTE_TIME        = 10    // Ticks after leaving ground where jump is still allowed
	INPUT_BUFFER_TICKS = 5     // Ticks before landing where a jump input is remembered
)

// PlayerMovementSystem handles all player movement mechanics including horizontal
// movement on flat ground and slopes, jumping with coyote time + early jump buffering,
// and platform drop-through functionality.
type PlayerMovementSystem struct{}

func (sys PlayerMovementSystem) Run(scene blueprint.Scene, dt float64) error {
	sys.handleHorizontal(scene)
	sys.handleJump(scene)
	return sys.handleDown(scene)
}

// handleHorizontal processes left/right movement with different behaviors for:
// - Air movement
// - Flat ground movement
// - Uphill/downhill slope movement with proper tangent calculations
func (PlayerMovementSystem) handleHorizontal(scene blueprint.Scene) {
	cursor := scene.NewCursor(blueprint.Queries.ActionBuffer)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		incomingActions := input.Components.ActionBuffer.GetFromCursor(cursor)
		direction := spatial.Components.Direction.GetFromCursor(cursor)

		_, pressedLeft := incomingActions.ConsumeAction(actions.Left)
		if pressedLeft {
			direction.SetLeft()
		}

		_, pressedRight := incomingActions.ConsumeAction(actions.Right)
		if pressedRight {
			direction.SetRight()
		}

		isMovingHorizontal := pressedLeft || pressedRight

		// Check ground status
		isGroundComponentPresent, onGround := components.OnGroundComponent.GetFromCursorSafe(cursor)
		isGrounded := isGroundComponentPresent && currentTick-1 == onGround.LastTouch

		// Horizontal airborne movement
		if !isGrounded {
			if isMovingHorizontal {
				dyn.Vel.X = SPEED_X * direction.AsFloat()
			}
			// Skip grounded horizontal movement
			continue
		}

		// Apply small downward force to keep player attached to slopes when grounded
		dyn.Vel.Y = math.Max(dyn.Vel.Y, SNAP_FORCE)

		// Horizontal flat movement
		flat := onGround.SlopeNormal.X == 0 && onGround.SlopeNormal.Y == 1
		if flat {
			if isMovingHorizontal {
				dyn.Vel.X = SPEED_X * direction.AsFloat()
			}
			// Skip slope horizontal movement
			continue
		}

		// Horizontal sloped movement
		if isMovingHorizontal {
			// Calculate tangent vector along the slope
			tangent := onGround.SlopeNormal.Perpendicular()

			isUphill := (direction.AsFloat() * onGround.SlopeNormal.X) > 0

			slopeDir := tangent.Scale(direction.AsFloat())

			if isUphill {
				// When going uphill, only set X velocity and let physics handle Y
				dyn.Vel.X = slopeDir.X * SPEED_X
			} else {
				// When going downhill, help player follow the slope with both X and Y velocities
				dyn.Vel.X = slopeDir.X * SPEED_X
				dyn.Vel.Y = slopeDir.Y * SPEED_X
			}
		}
	}
}

// handleJump processes jump inputs with coyote time and input buffering features
// Coyote time: Player can jump shortly after leaving a platform
// Input buffering: Jump inputs are remembered and applied when landing
func (PlayerMovementSystem) handleJump(scene blueprint.Scene) {
	playersEligibleToJumpQuery := warehouse.Factory.NewQuery().
		And(components.OnGroundComponent, input.Components.ActionBuffer)

	cursor := scene.NewCursor(playersEligibleToJumpQuery)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		// Get required components
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		incomingActions := input.Components.ActionBuffer.GetFromCursor(cursor)
		jumpState := components.JumpStateComponent.GetFromCursor(cursor)

		onGround := components.OnGroundComponent.GetFromCursor(cursor)

		if stampedAction, actionReceived := incomingActions.ConsumeAction(actions.Jump); actionReceived {

			// Coyote time: Allow jumping within certain ticks of leaving ground
			playerGroundedWithinCoyoteTime := currentTick-onGround.LastTouch <= COYOTE_TIME

			// Action buffering checks:
			//
			// 1. Was action received before touching ground?
			jumpIsBeforeGroundTouch := stampedAction.Tick <= onGround.LastTouch

			// 2. Was action within the buffer window?
			jumpWithinBufferWindow := onGround.LastTouch-stampedAction.Tick <= INPUT_BUFFER_TICKS

			// Combined buffer condition
			validBufferedJumpInput := jumpIsBeforeGroundTouch && jumpWithinBufferWindow

			// Direct jump: action received while already on ground
			directJumpAction := stampedAction.Tick >= onGround.LastTouch

			// Prevent double jumps: Make sure player hasn't jumped since last touching ground
			playerHasNotJumpedSinceGroundTouch := jumpState.LastJump < onGround.LastTouch

			// Player can jump if:
			// 1. They haven't jumped since touching ground, AND
			// 2a. They are in coyote time with a direct action, OR
			// 2b. They have a valid buffered action from before landing
			canJump := playerHasNotJumpedSinceGroundTouch &&
				((playerGroundedWithinCoyoteTime && directJumpAction) || validBufferedJumpInput)

			if canJump {
				dyn.Vel.Y = -JUMP_FORCE
				dyn.Accel.Y = -JUMP_FORCE
				jumpState.LastJump = currentTick
			}
		}
	}
}

// handleDown processes down input for platform drop-through functionality
// This allows players to press down to fall through one-way platforms
func (PlayerMovementSystem) handleDown(scene blueprint.Scene) error {
	playersEligibleToDropQuery := warehouse.Factory.NewQuery()
	playersEligibleToDropQuery.And(components.OnGroundComponent, input.Components.ActionBuffer)

	cursor := scene.NewCursor(playersEligibleToDropQuery)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		onGround := components.OnGroundComponent.GetFromCursor(cursor)
		jumpState := components.JumpStateComponent.GetFromCursor(cursor)

		// The component exists but were not actually grounded (will be removed soon)
		// Likely being held onto for coyote time tracking
		if onGround.LastTouch != currentTick-1 {
			continue
		}

		// Don't allow platform drop if player just jumped this tick
		if jumpState.LastJump == currentTick {
			continue
		}

		playerEntity, err := cursor.CurrentEntity()
		if err != nil {
			return err
		}

		incomingActions := input.Components.ActionBuffer.GetFromCursor(cursor)

		// Check for down action
		if stampedAction, inputReceived := incomingActions.ConsumeAction(actions.Down); inputReceived {
			// Only process actions fired within 5 ticks (slight tolerance for potential net delay)
			if currentTick-stampedAction.Tick > 5 {
				continue
			}

			// Otherwise, add the IgnorePlatform component to allow dropping through platforms
			err := playerEntity.EnqueueAddComponent(components.IgnorePlatformComponent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
