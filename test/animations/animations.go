package animations

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
)

var Up = client.AnimationData{
	Name:        "up",
	RowIndex:    3,
	FrameCount:  8,
	FrameWidth:  48,
	FrameHeight: 64,
	Speed:       8,
}

var Down = client.AnimationData{
	Name:        "down",
	RowIndex:    0,
	FrameCount:  8,
	FrameWidth:  48,
	FrameHeight: 64,
	Speed:       8,
}

var UpSide = client.AnimationData{
	Name:        "upside",
	RowIndex:    4,
	FrameCount:  8,
	FrameWidth:  48,
	FrameHeight: 64,
	Speed:       8,
}

var DownSide = client.AnimationData{
	Name:        "downside",
	RowIndex:    5,
	FrameCount:  8,
	FrameWidth:  48,
	FrameHeight: 64,
	Speed:       8,
}

// Note: downside and side are the same because I could only find a 6
// directional sprite. If we had an 8 directional sprite we could use a 4th
// animation here
var Side = client.AnimationData{
	Name:        "side",
	RowIndex:    5,
	FrameCount:  8,
	FrameWidth:  48,
	FrameHeight: 64,
	Speed:       8,
}
