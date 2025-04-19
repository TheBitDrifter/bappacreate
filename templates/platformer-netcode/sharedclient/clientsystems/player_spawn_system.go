package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/components"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/scenes"
)

type PlayerSpawnSystem struct{}

func (s PlayerSpawnSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	playerCursor := scene.NewCursor(blueprint.Queries.ActionBuffer)
	playerCount := playerCursor.TotalMatched()
	if playerCount != 0 {
		return nil
	}
	spawnQuery := warehouse.Factory.NewQuery().And(components.PlayerSpawnComponent)
	spawnCursor := scene.NewCursor(spawnQuery)

	var spawn components.PlayerSpawn
	for range spawnCursor.Next() {
		found := components.PlayerSpawnComponent.GetFromCursor(spawnCursor)
		spawn = *found
		break
	}

	_, err := scenes.NewPlayer(spawn.X, spawn.Y, scene.Storage())
	if err != nil {
		return err
	}
	return nil
}
