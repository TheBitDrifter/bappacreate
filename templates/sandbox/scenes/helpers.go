package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
)

// NewPlayer creates a 'player' entity for the scene
func NewPlayer(sto warehouse.Storage) error {
	playerArchetype, err := sto.NewOrExistingArchetype(
		ExamplePlayerComposition...,
	)
	err = playerArchetype.Generate(1,
		spatial.NewPosition(180, 180),
		spatial.NewRectangle(18, 58),
		motion.NewDynamics(10),
		spatial.NewDirectionRight(),
		input.InputBuffer{ReceiverIndex: 0},
		client.CameraIndex(0),
	)
	if err != nil {
		return err
	}
	return nil
}
