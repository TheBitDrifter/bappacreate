package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
)

var ExamplePlayerComposition = []warehouse.Component{
	spatial.Components.Position,
	client.Components.SpriteBundle,
	spatial.Components.Direction,
	input.Components.InputBuffer,
	client.Components.CameraIndex,
	spatial.Components.Shape,
	motion.Components.Dynamics,
	client.Components.SoundBundle,
}
