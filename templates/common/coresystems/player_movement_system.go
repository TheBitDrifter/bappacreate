package coresystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/common/actions"
	"github.com/TheBitDrifter/bappacreate/templates/common/components"
)

// PlayerMovementSystem handles all player movement mechanics including horizontal
// movement on flat ground and slopes, jumping with coyote time and input buffering,
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
	const (
		SPEED_X    = 120.0 // Player's horizontal movement speed
		SNAP_FORCE = 40.0  // Downward force to keep player attached to slopes
	)

	cursor := scene.NewCursor(blueprint.Queries.ActionBuffer)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		incomingInputs := input.Components.ActionBuffer.GetFromCursor(cursor)
		direction := spatial.Components.Direction.GetFromCursor(cursor)

		// Check and consume directional actions
		_, pressedLeft := incomingInputs.ConsumeAction(actions.Left)
		if pressedLeft {
			direction.SetLeft()
		}

		_, pressedRight := incomingInputs.ConsumeAction(actions.Right)
		if pressedRight {
			direction.SetRight()
		}

		// Track if player is attempting to move horizontally this frame
		isMovingHorizontal := pressedLeft || pressedRight

		// First check if the OnGroundComponent exists and get its value safely
		isGroundComponentPresent, onGround := components.OnGroundComponent.GetFromCursorSafe(cursor)
		isGrounded := isGroundComponentPresent && currentTick-1 == onGround.LastTouch

		// Handle Airborne:
		if !isGrounded {
			if isMovingHorizontal {
				dyn.Vel.X = SPEED_X * direction.AsFloat()
			}
			continue
		}

		// Handle Grounded:

		// Apply small downward force to keep player attached to slopes when grounded
		dyn.Vel.Y = math.Max(dyn.Vel.Y, SNAP_FORCE)

		// Handle flat
		flat := onGround.SlopeNormal.X == 0 && onGround.SlopeNormal.Y == 1
		if flat {
			if isMovingHorizontal {
				dyn.Vel.X = SPEED_X * direction.AsFloat()
			}
			continue
		}

		// Handle sloped
		if isMovingHorizontal {
			// Calculate tangent vector along the slope
			tangent := onGround.SlopeNormal.Perpendicular()

			// Determine if player is moving uphill by checking if direction and normal X have same sign
			isUphill := (direction.AsFloat() * onGround.SlopeNormal.X) > 0

			// Scale tangent by movement direction for correct slope alignment
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
func (PlayerMovementSystem) handleJump(scene blueprint.Scene) {
	const (
		JUMP_FORCE           = 320.0 // Upward force applied when jumping
		COYOTE_TIME_IN_TICKS = 10    // Ticks after leaving ground where jump is still allowed
		INPUT_BUFFER_TICKS   = 5     // Ticks before landing where a jump input is remembered
	)

	playersEligibleToJumpQuery := warehouse.Factory.NewQuery()
	playersEligibleToJumpQuery.And(components.OnGroundComponent, input.Components.ActionBuffer)

	cursor := scene.NewCursor(playersEligibleToJumpQuery)
	currentTick := scene.CurrentTick()

	for range cursor.Next() {
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		incomingInputs := input.Components.ActionBuffer.GetFromCursor(cursor)

		onGround := components.OnGroundComponent.GetFromCursor(cursor)

		if stampedInput, inputReceived := incomingInputs.ConsumeAction(actions.Jump); inputReceived {

			// Coyote time: Allow jumping within certain ticks of leaving ground
			playerGroundedWithinCoyoteTime := currentTick-onGround.LastTouch <= COYOTE_TIME_IN_TICKS

			// Input buffering checks:
			// 1. Was action received before touching ground?
			jumpInputIsBeforeGroundTouch := stampedInput.Tick <= onGround.LastTouch
			// 2. Was action within the buffer window?
			jumpInputWithinBufferWindow := onGround.LastTouch-stampedInput.Tick <= INPUT_BUFFER_TICKS

			validBufferedJumpInput := jumpInputIsBeforeGroundTouch && jumpInputWithinBufferWindow

			// Direct jump: Input received while already on ground
			directJumpInput := stampedInput.Tick >= onGround.LastTouch

			// Prevent double jumps: Make sure player hasn't jumped since last touching ground
			playerHasNotJumpedSinceGroundTouch := onGround.LastJump < onGround.LastTouch

			// Player can jump if:
			//
			// 1. They haven't jumped since touching ground, AND
			// 2a. They are in coyote time with a direct input, OR
			// 2b. They have a valid buffered input from before landing
			canJump := playerHasNotJumpedSinceGroundTouch &&
				((playerGroundedWithinCoyoteTime && directJumpInput) || validBufferedJumpInput)

			if canJump {
				dyn.Vel.Y = -JUMP_FORCE
				dyn.Accel.Y = -JUMP_FORCE
				onGround.LastJump = currentTick
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

		// The component exists but were not actually grounded (will be removed soon)
		// Likely being held onto for coyote time tracking
		if onGround.LastTouch != currentTick-1 {
			continue
		}

		// Don't allow platform drop if player just jumped this tick
		if onGround.LastJump == currentTick {
			continue
		}

		playerEntity, err := cursor.CurrentEntity()
		if err != nil {
			return err
		}

		incomingInputs := input.Components.ActionBuffer.GetFromCursor(cursor)

		// Check for down action
		if stampedInput, inputReceived := incomingInputs.ConsumeAction(actions.Down); inputReceived {
			// Only process inputs from this tick (ignore buffered inputs)
			if stampedInput.Tick != currentTick {
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
