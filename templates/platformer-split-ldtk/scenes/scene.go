package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/warehouse"

	"github.com/TheBitDrifter/bappa/blueprint/ldtk"
)

var entityRegistry = ldtk.NewLDtkEntityRegistry()

// Local scene object makes it easier to organize scene plans
type Scene struct {
	Name          string
	Plan          blueprint.Plan
	Width, Height int
}

// Registering custom LDTK entities
func init() {
	// Player start position handler
	entityRegistry.Register("PlayerStart", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		playerIndex := entity.IntFieldOr("playerIndex", 0) // Default to scene two if not specified
		return NewPlayer(float64(entity.Position[0]), float64(entity.Position[1]), playerIndex, sto)
	})

	// Scene transition trigger handler
	entityRegistry.Register("SceneTransfer", func(entity *ldtk.LDtkEntityInstance, sto warehouse.Storage) error {
		// Extract properties from LDtk entity
		targetScene := entity.StringFieldOr("targetScene", SCENE_TWO_NAME) // Default to scene two if not specified
		targetX := entity.FloatFieldOr("targetX", 20.0)
		targetY := entity.FloatFieldOr("targetY", 400.0)

		width := entity.FloatFieldOr("width", 100)
		height := entity.FloatFieldOr("height", 100)

		// Create the transfer trigger (scene change)
		return NewCollisionPlayerTransfer(
			sto,
			float64(entity.Position[0]),
			float64(entity.Position[1]),
			width,
			height,
			targetX,
			targetY,
			targetScene,
		)
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
