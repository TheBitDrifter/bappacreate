package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/warehouse"

	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/ldtk"
)

var entityRegistry = ldtk.NewLDtkEntityRegistry()

// Local scene object makes it easier to organize scene plans
type Scene struct {
	Name          string
	Plan          blueprint.Plan
	Width, Height int
	// Optional: manual assets to preload (usually for entities that get added 'dynamically', while the scene is running)
	Preload client.PreLoadAssetBundle
}

// Registering custom LDTK entities
func init() {
	// Player start position handler
	entityRegistry.Register("PlayerStart", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		// Create the player at the position defined in LDtk
		_, err := NewPlayerSpawn(float64(entity.Position[0]), float64(entity.Position[1]), sto)
		if err != nil {
			return err
		}
		return nil
	})

	// Ramp
	entityRegistry.Register("Ramp", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		return NewRamp(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
		)
	})

	// RotatedPlatform
	entityRegistry.Register("RotatedPlatform", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		return NewPlatformRotated(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			-0.25,
		)
	})
}
