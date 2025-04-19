package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/topdown/components"
	"github.com/TheBitDrifter/bappacreate/templates/topdown/sounds"
)

type PlayerSoundSystem struct{}

func (sys PlayerSoundSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	const WALK_VOL = 0.6

	playersMovingWithSounds := warehouse.Factory.NewQuery().And(
		client.Components.SoundBundle,
		components.IsMovingComponent,
	)

	cursor := scene.NewCursor(playersMovingWithSounds)

	for range cursor.Next() {
		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)

		runSound, err := coldbrew.MaterializeSound(soundBundle, sounds.Run)
		if err != nil {
			return err
		}

		audPlayer := runSound.GetAny()

		// Play the sound
		if !audPlayer.IsPlaying() {
			audPlayer.SetVolume(WALK_VOL)
			audPlayer.Rewind()
			audPlayer.Play()
		}
	}
	return nil
}
