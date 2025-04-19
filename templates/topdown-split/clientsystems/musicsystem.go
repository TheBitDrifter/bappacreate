package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/components"
	"github.com/TheBitDrifter/bappacreate/templates/topdown-split/sounds"
)

type MusicSystem struct{}

// Note: this a very simple music system that does not account for multiple scenes
// Adjust accordingly
func (sys MusicSystem) Run(lc coldbrew.LocalClient, scene coldbrew.Scene) error {
	musicQuery := warehouse.Factory.NewQuery().And(components.MusicTag)
	cursor := scene.NewCursor(musicQuery)

	// There's only one but iterate nonetheless
	for range cursor.Next() {
		soundBundle := client.Components.SoundBundle.GetFromCursor(cursor)

		sound, err := coldbrew.MaterializeSound(soundBundle, sounds.Music)
		if err != nil {
			return err
		}
		player := sound.GetAny()

		// Loop if needed
		if !player.IsPlaying() {
			player.Rewind()
			player.Play()
		}
	}
	return nil
}
