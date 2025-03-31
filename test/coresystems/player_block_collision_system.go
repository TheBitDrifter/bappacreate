package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/test/components"
)

type PlayerBlockCollisionSystem struct{}

func (s PlayerBlockCollisionSystem) Run(scene blueprint.Scene, dt float64) error {
	// Create cursors
	blockTerrainQuery := warehouse.Factory.NewQuery().And(components.BlockTerrainTag)
	blockTerrainCursor := scene.NewCursor(blockTerrainQuery)
	playerCursor := scene.NewCursor(blueprint.Queries.InputBuffer)

	// Outer loop is blocks
	for range blockTerrainCursor.Next() {
		// Inner is players
		for range playerCursor.Next() {
			// Delegate to helper
			err := s.resolve(scene, blockTerrainCursor, playerCursor)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Main collision logic
func (PlayerBlockCollisionSystem) resolve(scene blueprint.Scene, blockCursor, playerCursor *warehouse.Cursor) error {
	// Get the player pos, shape, and dynamics
	playerPosition := spatial.Components.Position.GetFromCursor(playerCursor)
	playerShape := spatial.Components.Shape.GetFromCursor(playerCursor)
	playerDynamics := motion.Components.Dynamics.GetFromCursor(playerCursor)

	// Get the block pos, shape, and dynamics
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
