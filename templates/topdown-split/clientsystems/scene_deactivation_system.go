package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
)

type SceneDeactivationSystem struct{}

// Global client system deactivates scenes with no players
// Watch out for this disabling scenes that are UI based only
// You may need to handle that (add a hidden player to the scene or modify the check)!
func (SceneDeactivationSystem) Run(cli coldbrew.Client) error {
	for scene := range cli.ActiveScenes() {
		cursor := warehouse.Factory.NewCursor(blueprint.Queries.ActionBuffer, scene.Storage())
		hasPlayers := cursor.TotalMatched() > 0
		if !hasPlayers {
			cli.DeactivateScene(scene)
		}
	}
	return nil
}
