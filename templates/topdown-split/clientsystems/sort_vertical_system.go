package clientsystems

import (
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/tteokbokki/spatial"
	"github.com/TheBitDrifter/bappa/warehouse"
)

type SortVerticalSystem struct{}

// This system assigns priority to sprites via their vertical position (y)
// This allows the GlobalRenderSystem to draw them in proper order
func (SortVerticalSystem) Run(cli coldbrew.LocalClient, scene coldbrew.Scene) error {
	query := warehouse.Factory.NewQuery().And(
		client.Components.SpriteBundle,
		spatial.Components.Position,
		warehouse.Factory.NewQuery().Not(client.Components.ParallaxBackground),
	)

	cursor := scene.NewCursor(query)

	for range cursor.Next() {
		pos := spatial.Components.Position.GetFromCursor(cursor)
		bundle := client.Components.SpriteBundle.GetFromCursor(cursor)

		// Assign priority based on posY
		priority := pos.Y
		for i := range bundle.Blueprints {
			bp := &bundle.Blueprints[i]
			bp.Config.Priority = int(priority) + scene.Height()
		}
	}
	return nil
}
