package components

import (
	"github.com/TheBitDrifter/bappa/warehouse"
)

var ExampleComponent = warehouse.FactoryNewComponent[Example]()

type Example struct {
	SomeState int
}
