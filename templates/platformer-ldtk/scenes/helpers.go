package scenes

import (
	"github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/animations"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/components"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/sounds"
	"github.com/TheBitDrifter/blueprint"
	blueprintclient "github.com/TheBitDrifter/blueprint/client"
	blueprintinput "github.com/TheBitDrifter/blueprint/input"
	blueprintmotion "github.com/TheBitDrifter/blueprint/motion"
	blueprintspatial "github.com/TheBitDrifter/blueprint/spatial"
	"github.com/TheBitDrifter/blueprint/vector"
	"github.com/TheBitDrifter/warehouse"
)

// NewPlayer creates a player entity for the scene
func NewPlayer(x, y float64, sto warehouse.Storage) error {
	playerArchetype, err := sto.NewOrExistingArchetype(
		PlayerComposition...,
	)
	err = playerArchetype.Generate(1,
		blueprintspatial.NewPosition(x, y),
		blueprintspatial.NewRectangle(18, 58),
		blueprintmotion.NewDynamics(10),
		blueprintspatial.NewDirectionRight(),
		blueprintinput.InputBuffer{ReceiverIndex: 0},
		blueprintclient.CameraIndex(0),
		blueprintclient.NewSpriteBundle().
			AddSprite("characters/box_man_sheet.png", true).
			WithAnimations(animations.IdleAnimation, animations.RunAnimation, animations.FallAnimation, animations.JumpAnimation).
			SetActiveAnimation(animations.IdleAnimation).
			WithOffset(vector.Two{X: -72, Y: -59}).
			WithPriority(20),
		blueprintclient.NewSoundBundle().
			AddSoundFromConfig(sounds.Run).
			AddSoundFromConfig(sounds.Jump).
			AddSoundFromConfig(sounds.Land),
	)
	if err != nil {
		return err
	}
	return nil
}

// NewPlatformRotated creates a one way platform
func NewPlatformRotated(sto warehouse.Storage, x, y, rotation float64) error {
	platformArche, err := sto.NewOrExistingArchetype(PlatformComposition...)
	if err != nil {
		return err
	}
	return platformArche.Generate(1,
		blueprintspatial.NewPosition(x, y),
		blueprintspatial.Rotation(rotation),
		// Triangles for one way platform
		blueprintspatial.NewTriangularPlatform(144, 16),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/platform.png", true).
			WithOffset(vector.Two{X: -72, Y: -8}),
	)
}

// NewRamp creates a ramp (sloped block hexagon)
func NewRamp(sto warehouse.Storage, x, y float64) error {
	// Add a sprite
	composition := []warehouse.Component{
		blueprintclient.Components.SpriteBundle,
	}

	// Compose the archetype with the sprite and block composition
	composition = append(composition, BlockTerrainComposition...)
	rampArche, err := sto.NewOrExistingArchetype(composition...)
	if err != nil {
		return err
	}

	return rampArche.Generate(1,
		blueprintspatial.NewPosition(x, y),
		blueprintspatial.NewDoubleRamp(250, 46, 0.2),
		blueprintclient.NewSpriteBundle().
			AddSprite("terrain/ramp.png", true).
			WithOffset(vector.Two{X: -125, Y: -22}),
	)
}

// NewCityBackground creates the city parallax background entities
func NewCityBackground(sto warehouse.Storage) error {
	return blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("backgrounds/city/sky.png", 0.025, 0.025).
		AddLayer("backgrounds/city/far.png", 0.025, 0.05).
		AddLayer("backgrounds/city/mid.png", 0.1, 0.1).
		AddLayer("backgrounds/city/near.png", 0.2, 0.2).
		Build()
}

// NewSkyBackground creates a sky background entity
func NewSkyBackground(sto warehouse.Storage) error {
	return blueprint.NewParallaxBackgroundBuilder(sto).
		AddLayer("backgrounds/city/sky.png", 0.05, 0.05).
		Build()
}

// NewJazzMusic adds a Jazz music entity
func NewJazzMusic(sto warehouse.Storage) error {
	musicArche, err := sto.NewOrExistingArchetype(MusicComposition...)
	if err != nil {
		return err
	}
	return musicArche.Generate(1, blueprintclient.NewSoundBundle().AddSoundFromPath("music.wav"))
}

// NewCollisionPlayerTransfer creates an collidable entity/shape that will transfer the player
// to the targeted pos and scene upon touching it
func NewCollisionPlayerTransfer(
	sto warehouse.Storage, x, y, w, h, playerTargetX, playerTargetY float64, target string,
) error {
	collisionPlayerTransferArche, err := sto.NewOrExistingArchetype(
		CollisionPlayerTransferComposition...,
	)
	if err != nil {
		return err
	}
	return collisionPlayerTransferArche.Generate(1,
		blueprintspatial.NewPosition(x, y),
		blueprintspatial.NewRectangle(w, h),
		components.PlayerSceneTransfer{
			Dest: target,
			X:    playerTargetX,
			Y:    playerTargetY,
		},
	)
}
