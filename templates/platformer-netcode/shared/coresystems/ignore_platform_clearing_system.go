package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/components"
)

// IgnorePlatformClearingSystem clears out individual platform ignores that have expired
type IgnorePlatformClearingSystem struct{}

func (IgnorePlatformClearingSystem) Run(scene blueprint.Scene, dt float64) error {
	ignorePlatformQuery := warehouse.Factory.NewQuery().And(components.IgnorePlatformComponent)
	ignorePlatformCursor := scene.NewCursor(ignorePlatformQuery)

	const expirationTicks = 15

	for range ignorePlatformCursor.Next() {
		ignorePlatform := components.IgnorePlatformComponent.GetFromCursor(ignorePlatformCursor)
		currentTick := scene.CurrentTick()

		// Track if we have any active ignores left
		anyActive := false

		for i := range ignorePlatform.Items {
			// Skip already cleared entries
			if ignorePlatform.Items[i].EntityID == 0 {
				continue
			}

			// Check if this entry has expired
			if currentTick-ignorePlatform.Items[i].LastActive > expirationTicks {

				// Clear this specific entry
				ignorePlatform.Items[i].EntityID = 0
				ignorePlatform.Items[i].Recycled = 0
				ignorePlatform.Items[i].LastActive = 0

			} else {
				anyActive = true
			}
		}
		// If we don't have any active ignores left, remove the entire component
		if !anyActive {
			ignoringEntity, _ := ignorePlatformCursor.CurrentEntity()
			err := ignoringEntity.EnqueueRemoveComponent(components.IgnorePlatformComponent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
