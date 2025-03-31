package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/tteo_coresystems"
)

var DefaultCoreSystems = []blueprint.CoreSystem{
	PlayerMovementSystem{},               // Apply player input forces
	tteo_coresystems.IntegrationSystem{}, // Update velocities and positions
	tteo_coresystems.TransformSystem{},   // Update collision shapes
	PlayerBlockCollisionSystem{},         // Handle  collisions
}
