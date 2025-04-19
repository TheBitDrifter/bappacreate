package components

import (
	"github.com/TheBitDrifter/bappa/warehouse"
)

var (
	PlayerSceneTransferComponent = warehouse.FactoryNewComponent[PlayerSceneTransfer]()
	DirectionEightComponent      = warehouse.FactoryNewComponent[DirectionEight]()
	IsMovingComponent            = warehouse.FactoryNewComponent[struct{}]()
)
