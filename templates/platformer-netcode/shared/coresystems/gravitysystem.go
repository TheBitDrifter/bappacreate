package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/tteokbokki/motion"
)

const (
	DEFAULT_GRAVITY  = 9.8
	PIXELS_PER_METER = 50.0
)

type GravitySystem struct{}

func (GravitySystem) Run(scene blueprint.Scene, dt float64) error {
	cursor := scene.NewCursor(blueprint.Queries.Dynamics)

	for range cursor.Next() {
		dyn := motion.Components.Dynamics.GetFromCursor(cursor)

		mass := 1 / dyn.InverseMass

		gravity := motion.Forces.Generator.NewGravityForce(mass, DEFAULT_GRAVITY, PIXELS_PER_METER)

		motion.Forces.AddForce(dyn, gravity)
	}
	return nil
}
