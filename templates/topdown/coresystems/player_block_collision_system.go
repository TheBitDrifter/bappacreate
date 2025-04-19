package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/topdown/components"
)

type PlayerBlockCollisionSystem struct{}

func (s PlayerBlockCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	blockTerrainQuery := warehouse.Factory.NewQuery().And(components.BlockTerrainTag)
	blockTerrainCursor := scene.NewCursor(blockTerrainQuery)
	playerCursor := scene.NewCursor(blueprint.Queries.ActionBuffer)

	for range blockTerrainCursor.Next() {
		for range playerCursor.Next() {
			err := s.resolve(scene, blockTerrainCursor, playerCursor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (PlayerBlockCollisionSystem) resolve(scene blueprint.Scene, blockCursor, playerCursor *warehouse.Cursor) error {
	playerPosition := spatial.Components.Position.GetFromCursor(playerCursor)
	playerShape := spatial.Components.Shape.GetFromCursor(playerCursor)
	playerDynamics := motion.Components.Dynamics.GetFromCursor(playerCursor)

	blockPosition := spatial.Components.Position.GetFromCursor(blockCursor)
	blockShape := spatial.Components.Shape.GetFromCursor(blockCursor)
	blockDynamics := motion.Components.Dynamics.GetFromCursor(blockCursor)

	// Check for a collision
	if ok, collisionResult := spatial.Detector.Check(
		*playerShape, *blockShape, playerPosition.Two, blockPosition.Two,
	); ok {
		// Resolve a collision
		motion.Resolver.Resolve(
			&playerPosition.Two,
			&blockPosition.Two,
			playerDynamics,
			blockDynamics,
			collisionResult,
		)
	}
	return nil
}
