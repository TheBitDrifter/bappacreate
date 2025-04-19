package clientsystems

import (
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_clientsystems"
)

var DefaultClientSystems = []coldbrew.ClientSystem{
	PlayerSoundSystem{},
	MusicSystem{},
	PlayerAnimationSystem{},
	&CameraFollowerSystem{},
	&coldbrew_clientsystems.BackgroundScrollSystem{},
	PlayerSpawnSystem{},
}

var DefaultClientSystemsNetworked = []coldbrew.ClientSystem{
	PlayerSoundSystem{},
	MusicSystem{},
	PlayerAnimationSystem{},
	&CameraFollowerSystem{},
	&coldbrew_clientsystems.BackgroundScrollSystem{},
}
