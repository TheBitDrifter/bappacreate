package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/tteo_coresystems"
)

var DefaultCoreSystems = []blueprint.CoreSystem{
	// If you want default physics:
	tteo_coresystems.IntegrationSystem{},
	tteo_coresystems.TransformSystem{},
}
