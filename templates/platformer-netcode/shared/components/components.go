package components

import (
	"github.com/TheBitDrifter/bappa/warehouse"
)

var (
	OnGroundComponent            = warehouse.FactoryNewComponent[OnGround]()
	IgnorePlatformComponent      = warehouse.FactoryNewComponent[IgnorePlatform]()
	PlayerSceneTransferComponent = warehouse.FactoryNewComponent[PlayerSceneTransfer]()
	JumpStateComponent           = warehouse.FactoryNewComponent[JumpState]()
	PlayerSpawnComponent         = warehouse.FactoryNewComponent[PlayerSpawn]()
)
