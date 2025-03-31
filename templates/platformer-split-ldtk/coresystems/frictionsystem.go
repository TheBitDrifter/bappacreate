package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
)

const (
	DEFAULT_FRICTION = 0.5
	DEFAULT_DAMP     = 0.9
)

type FrictionSystem struct{}

func (FrictionSystem) Run(scene blueprint.Scene, dt float64) error {
	// Iterate through entities with dynamics components(physics)
	cursor := scene.NewCursor(blueprint.Queries.Dynamics)
	for range cursor.Next() {
		// Get the dynamics
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)
		friction := motion.Forces.Generator.NewHorizontalFrictionForce(dyn.Vel, DEFAULT_FRICTION)
		motion.Forces.AddForce(dyn, friction)

		motion.Forces.Generator.ApplyHorizontalDamping(dyn, DEFAULT_DAMP)
	}
	return nil
}
