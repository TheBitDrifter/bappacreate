package rendersystems

import "github.com/TheBitDrifter/bappa/coldbrew"

var DefaultRenderSystems = []coldbrew.RenderSystem{
	PlayerCameraPriorityRenderer{},
}
