package rendersystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/blueprint/vector"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/coldbrew/coldbrew_rendersystems"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
)

type PlayerCameraPriorityRenderer struct{}

func (PlayerCameraPriorityRenderer) Render(scene coldbrew.Scene, screen coldbrew.Screen, c coldbrew.LocalClient) {
	netCli, isNet := c.(coldbrew.NetworkClient)
	if isNet {
		id, ok := netCli.AssociatedEntityID()
		if !ok {
			return
		}

		pEn, err := scene.Storage().Entity(id)
		if err != nil {
			return
		}
		for _, cam := range netCli.ActiveCamerasFor(scene) {
			// If it ain't ready chill out!
			if !netCli.Ready(cam) {
				continue
			}
			bundle := client.Components.SpriteBundle.GetFromEntity(pEn)
			spr := coldbrew.MaterializeSprites(bundle)[0]

			coldbrew_rendersystems.RenderEntity(
				spatial.Components.Position.GetFromEntity(pEn).Two,
				0,
				vector.Two{1, 1},
				*spatial.Components.Direction.GetFromEntity(pEn),
				spr,
				&bundle.Blueprints[0],
				cam,
				scene.CurrentTick(),
			)
			cam.PresentToScreen(screen, coldbrew.ClientConfig.CameraBorderSize())
		}

	}
}
