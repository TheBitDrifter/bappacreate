package scenes

import (
	"log"

	"github.com/TheBitDrifter/bappa/warehouse"
	"github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/ldtk"
)

const SCENE_TWO_NAME = "Scene2"

var SceneTwo = Scene{
	Name:   SCENE_TWO_NAME,
	Plan:   sceneTwoPlan,
	Width:  ldtk.DATA.WidthFor(SCENE_TWO_NAME),
	Height: ldtk.DATA.HeightFor(SCENE_TWO_NAME),
}

// Scene two is a simple night sky and floor
func sceneTwoPlan(width, height int, sto warehouse.Storage) error {
	// Load the image tiles
	err := ldtk.DATA.LoadTiles(SCENE_TWO_NAME, sto)
	if err != nil {
		return err
	}

	// Load the terrain
	// Pass the terrain archetypes in order of int grid layer they map to
	blockArchetype, _ := sto.NewOrExistingArchetype(BlockTerrainComposition...)
	platArchetype, _ := sto.NewOrExistingArchetype(PlatformComposition...)
	transferArchetype, _ := sto.NewOrExistingArchetype(CollisionPlayerTransferComposition...)

	err = ldtk.DATA.LoadIntGrid(SCENE_TWO_NAME, sto, blockArchetype, platArchetype, transferArchetype)
	if err != nil {
		return err
	}

	// Load custom LDTK entities
	err = ldtk.DATA.LoadEntities(SCENE_TWO_NAME, sto, entityRegistry)
	if err != nil {
		log.Printf("Error loading entities: %v", err)
		return err
	}

	// Background
	return NewSkyBackground(sto)
}
