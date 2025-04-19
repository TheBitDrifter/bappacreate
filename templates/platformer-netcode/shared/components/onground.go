package components

import "github.com/TheBitDrifter/bappa/blueprint/vector"

type OnGround struct {
	LastTouch   int
	Landed      int
	SlopeNormal vector.Two
}
