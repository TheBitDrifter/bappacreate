package main

import (
	"encoding/json"
	"log"

	"github.com/TheBitDrifter/bappa/blueprint"
	"github.com/TheBitDrifter/bappa/blueprint/client"
	"github.com/TheBitDrifter/bappa/coldbrew"
	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared/scenes"
)

func Derser(nc coldbrew.NetworkClient, data []byte) error {
	activeScenes := nc.ActiveScenes()
	var scene coldbrew.Scene
	for s := range activeScenes {
		scene = s
		break
	}
	if scene != nil && scene.Ready() {
		storage := scene.Storage()
		if storage != nil {
			var world warehouse.SerializedStorage
			err := json.Unmarshal(data, &world)
			if err != nil {
				log.Printf("NetworkClient Update Error: Failed to unmarshal state (%d bytes): %v", len(data), err)
			} else {

				seen := map[int]struct{}{}

				for _, se := range world.Entities {
					seen[int(se.ID)] = struct{}{}

					en, err := storage.ForceSerializedEntityExclude(
						se, client.Components.SoundBundle,
						client.Components.SpriteBundle,
					)
					if err != nil {
						return err
					}

					err = se.SetValue(en)
					if err != nil {
						return err
					}

					if !en.Table().Contains(client.Components.SpriteBundle) {
						err := en.AddComponentWithValue(client.Components.SpriteBundle, scenes.DEFAULT_PLAYER_SPR_BUNDLE)
						if err != nil {
							return err
						}

						err = en.AddComponentWithValue(client.Components.SoundBundle, scenes.DEFAULT_PLAYER_SND_BUNDLE)
						if err != nil {
							return err
						}
					}

				}

				purge := []warehouse.Entity{}
				query := blueprint.Queries.ActionBuffer
				cursor := scene.NewCursor(query)

				for range cursor.Next() {
					e, _ := cursor.CurrentEntity()
					if _, ok := seen[int(e.ID())]; !ok {
						purge = append(purge, e)
					}
				}

				err := storage.DestroyEntities(purge...)
				if err != nil {
					log.Println(err)
				}

				coldbrew.ForceSetTick(world.CurrentTick)
			}
		} else {
			log.Println("NetworkClient Update Error: Active scene has nil storage.")
		}
	}
	return nil
}
