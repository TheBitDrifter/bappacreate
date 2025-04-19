package components

import "github.com/TheBitDrifter/bappa/warehouse"

// Tags help us identify/categorize archetypes/entities when their
// composition alone isn't enough.
//
// For example its hard to tell the
// difference between a block and platform since they both have
// dynamics, shapes, sprites, etc
type blockTag struct{}

type platTag struct{}

type musicTag struct{}

var (
	BlockTerrainTag = warehouse.FactoryNewComponent[blockTag]()
	PlatformTag     = warehouse.FactoryNewComponent[platTag]()
	MusicTag        = warehouse.FactoryNewComponent[musicTag]()
)
