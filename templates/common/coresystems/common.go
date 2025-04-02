package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/tteo_coresystems"
)

var DefaultCoreSystems = []blueprint.CoreSystem{
	GravitySystem{},                      // Apply gravity forces
	FrictionSystem{},                     // Apply Friction forces
	PlayerMovementSystem{},               // Apply player input forces
	tteo_coresystems.IntegrationSystem{}, // Update velocities and positions
	tteo_coresystems.TransformSystem{},   // Update collision shapes
	PlayerBlockCollisionSystem{},         // Handle  collisions
	NewPlayerPlatformCollisionSystem(),   // Handle  collisions — func returns ptr because system is not pure (has state)
	OnGroundClearingSystem{},             // Clear onGround
	IgnorePlatformClearingSystem{},       // Clear ignorePlatform
}
