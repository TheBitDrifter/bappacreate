package clientsystems

import (
	"math"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/common/components"
	"github.com/TheBitDrifter/bappacreate/templates/common/sounds"
)

type PlayerSoundSystem struct{}

func (sys PlayerSoundSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	// Create a query for players that can play sounds and are on the ground
	playersWithSoundsOnTheGround := warehouse.Factory.NewQuery()
	playersWithSoundsOnTheGround.And(
		client.Components.SoundBundle, // Has sounds
		input.Components.InputBuffer,  // Can receive input
		motion.Components.Dynamics,    // Has physics properties
		components.OnGroundComponent,  // Is on the ground
	)

	// Get all entities that match the query
	cursor := scene.NewCursor(playersWithSoundsOnTheGround)

	// Iterate over all matching players
	for range cursor.Next() {

		// Get state
		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		onGround := components.OnGroundComponent.GetFromCursor(cursor)
		currentTick := scene.CurrentTick()

		// Landed sound
		if onGround.Landed == currentTick {
			landingSound, _ := coldbrew.MaterializeSound(soundBundle, sounds.Land)
			player := landingSound.GetAny()

			// A hack to prevent landing sound artifacts between scenes
			// In a more robust setup, we might track if a player has recently changed scenes via a component
			// Such a component would be helpful here
			sceneRecentlySelected := scene.CurrentTick()-scene.LastSelectedTick() < 30

			if !player.IsPlaying() && !sceneRecentlySelected {
				player.Rewind()
				player.Play()
			}
		}

		// Jump sound

		// Gotta have y velocity
		if dyn.Vel.Y < 5 && onGround.LastJump == currentTick {
			jumpSound, _ := coldbrew.MaterializeSound(soundBundle, sounds.Jump)
			player := jumpSound.GetAny()

			if !player.IsPlaying() {
				player.Rewind()
				player.Play()
			}
		}

		// Run Sound
		// Must be moving horizontally
		const minMovementSpeed = 20.0
		if math.Abs(dyn.Vel.X) <= minMovementSpeed {
			continue
		}
		// Ensure onGround is not just available for coyote timer
		touchedGroundThisTick := onGround.LastTouch == currentTick
		if !touchedGroundThisTick {
			continue
		}

		runSound, err := coldbrew.MaterializeSound(soundBundle, sounds.Run)
		if err != nil {
			return err
		}
		player := runSound.GetAny()

		if !player.IsPlaying() {
			player.Rewind()
			player.Play()
		}
	}

	return nil
}
