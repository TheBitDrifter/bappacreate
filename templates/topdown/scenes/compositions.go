package scenes

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/topdown/components"
)

// These are slices of common component compositions for various archetypes.
// They only include/represent the initial and static components of archetype
// Components can still be added or removed dynamically at runtime
//
// These slices are especially useful for creating starting entities, via archetypes, inside plan functions

var PlayerComposition = []warehouse.Component{
	spatial.Components.Position,
	motion.Components.Dynamics,
	client.Components.SpriteBundle,
	input.Components.ActionBuffer,
	client.Components.CameraIndex,
	spatial.Components.Shape,
	client.Components.SoundBundle,
	spatial.Components.Direction,
	components.DirectionEightComponent,
}

var BlockTerrainComposition = []warehouse.Component{
	components.BlockTerrainTag,
	spatial.Components.Shape,
	spatial.Components.Position,
	motion.Components.Dynamics,
}

var PropComposition = []warehouse.Component{
	components.BlockTerrainTag,
	client.Components.SpriteBundle,
	spatial.Components.Shape,
	spatial.Components.Position,
	motion.Components.Dynamics,
}

var CollisionPlayerTransferComposition = []warehouse.Component{
	spatial.Components.Position,
	spatial.Components.Shape,
	components.PlayerSceneTransferComponent,
}

var MusicComposition = []warehouse.Component{
	client.Components.SoundBundle,
	components.MusicTag,
}
