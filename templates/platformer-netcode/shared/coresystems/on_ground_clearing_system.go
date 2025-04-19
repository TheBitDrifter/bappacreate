package coresystems

import (
	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/components"
)

type OnGroundClearingSystem struct{}

func (OnGroundClearingSystem) Run(scene blueprint.Scene, dt float64) error {
	const EXPIRATION_IN_TICKS = 15

	onGroundQuery := warehouse.Factory.NewQuery().And(components.OnGroundComponent)
	onGroundCursor := scene.NewCursor(onGroundQuery)

	for range onGroundCursor.Next() {
		onGround := components.OnGroundComponent.GetFromCursor(onGroundCursor)

		// If it's expired, remove it
		if scene.CurrentTick()-onGround.LastTouch > EXPIRATION_IN_TICKS {
			groundedEntity, _ := onGroundCursor.CurrentEntity()

			err := groundedEntity.EnqueueRemoveComponent(components.OnGroundComponent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
