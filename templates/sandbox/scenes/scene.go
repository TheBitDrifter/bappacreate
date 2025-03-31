package scenes

import "github.com/TheBitDrifter/bappa/blueprint"

// Local scene object makes it easier to organize scene plans
type Scene struct {
	Name          string
	Plan          blueprint.Plan
	Width, Height int
}
