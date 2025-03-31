package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/components"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/sounds"
)

type PlayerSoundSystem struct{}

func (sys PlayerSoundSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	const WALK_VOL = 0.6

	// Create a query for players that can play sounds and are moving
	playersMovingWithSounds := warehouse.Factory.NewQuery()
	playersMovingWithSounds.And(
		client.Components.SoundBundle, // Has sounds
		components.IsMovingComponent,  // Is moving
	)

	// Get all entities that match the query
	cursor := scene.NewCursor(playersMovingWithSounds)

	// Iterate
	for range cursor.Next() {
		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)

		// Get the sound from the bundle
		runSound, err := coldbrew.MaterializeSound(soundBundle, sounds.Run)
		if err != nil {
			return err
		}

		// Get a player from the sound
		player := runSound.GetAny()

		// Play the sound
		if !player.IsPlaying() {
			player.SetVolume(WALK_VOL)
			player.Rewind()
			player.Play()
		}
	}
	return nil
}
