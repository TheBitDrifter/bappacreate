package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/topdown/animations"
	"github.com/TheBitDrifter/bappacreate/templates/topdown/components"
	"github.com/TheBitDrifter/bappacreate/templates/topdown/sounds"
)

// NewPlayer creates a player entity
func NewPlayer(sto warehouse.Storage) error {
	playerArchetype, err := sto.NewOrExistingArchetype(
		PlayerComposition...,
	)
	err = playerArchetype.Generate(1,
		spatial.NewPosition(180, 180),
		spatial.NewRectangle(16, 16),
		motion.NewDynamics(10),
		spatial.NewDirectionRight(),
		input.InputBuffer{ReceiverIndex: 0},
		client.CameraIndex(0),
		client.NewSpriteBundle().
			AddSprite("characters/main/idle.png", true).
			WithAnimations(animations.Down, animations.Side, animations.DownSide, animations.UpSide, animations.Up).
			SetActiveAnimation(animations.Down).
			WithOffset(vector.Two{X: -24, Y: -32}).
			AddSprite("characters/main/walk.png", false).
			WithAnimations(animations.Down, animations.Side, animations.DownSide, animations.UpSide, animations.Up).
			SetActiveAnimation(animations.Down).
			WithOffset(vector.Two{X: -24, Y: -32}),

		client.NewSoundBundle().AddSoundFromConfig(sounds.Run),

		components.NewDirectionDown(),
	)
	if err != nil {
		return err
	}
	return nil
}

// NewTreeProp creates a tree prop entity
func NewTreeProp(sto warehouse.Storage, x, y float64) error {
	treePropArche, err := sto.NewOrExistingArchetype(
		PropComposition...,
	)
	if err != nil {
		return err
	}
	return treePropArche.Generate(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(10, 10),
		motion.NewDynamics(0),
		client.NewSpriteBundle().
			AddSprite("props/tree.png", true).
			WithOffset(vector.Two{X: -45, Y: -130}),
	)
}

// NewMoveableStatueProp creates a moveable statue prop entity
func NewMoveableStatueProp(sto warehouse.Storage, x, y float64) error {
	statueArche, err := sto.NewOrExistingArchetype(
		PropComposition...,
	)
	if err != nil {
		return err
	}
	return statueArche.Generate(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(28, 20),
		motion.NewDynamics(10),
		client.NewSpriteBundle().
			AddSprite("props/statue.png", true).
			WithOffset(vector.Two{X: -17, Y: -60}),
	)
}

// NewBlockTerrain creates and invisible bounds block entity
func NewBlockTerrain(sto warehouse.Storage, x, y, w, h float64) error {
	statueArche, err := sto.NewOrExistingArchetype(
		BlockTerrainComposition...,
	)
	if err != nil {
		return err
	}
	return statueArche.Generate(1,
		spatial.NewPosition(x, y),
		spatial.NewRectangle(w, h),
		motion.NewDynamics(0),
	)
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
		spatial.NewPosition(x, y),
		spatial.NewRectangle(w, h),
		components.PlayerSceneTransfer{
			Dest: target,
			X:    playerTargetX,
			Y:    playerTargetY,
		},
	)
}

// NewFantasyMusic creates a simple music entity with the default fantasy music wav
func NewFantasyMusic(sto warehouse.Storage) error {
	musicArche, err := sto.NewOrExistingArchetype(MusicComposition...)
	if err != nil {
		return err
	}
	return musicArche.Generate(1, client.NewSoundBundle().AddSoundFromConfig(sounds.Music))
}
