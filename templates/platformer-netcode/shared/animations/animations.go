package animations

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
)

var IdleAnimation = client.AnimationData{
	Name:        "idle",
	RowIndex:    0,
	FrameCount:  6,
	FrameWidth:  144,
	FrameHeight: 116,
	Speed:       8,
}

var RunAnimation = client.AnimationData{
	Name:        "run",
	RowIndex:    1,
	FrameCount:  8,
	FrameWidth:  144,
	FrameHeight: 116,
	Speed:       5,
}

var JumpAnimation = client.AnimationData{
	Name:           "jump",
	RowIndex:       2,
	FrameCount:     3,
	FrameWidth:     144,
	FrameHeight:    116,
	Speed:          5,
	Freeze:         true,
	PositionOffset: vector.Two{X: 0, Y: 10},
}

var FallAnimation = client.AnimationData{
	Name:           "fall",
	RowIndex:       3,
	FrameCount:     3,
	FrameWidth:     144,
	FrameHeight:    116,
	Speed:          5,
	Freeze:         true,
	PositionOffset: vector.Two{X: 0, Y: 10},
}
