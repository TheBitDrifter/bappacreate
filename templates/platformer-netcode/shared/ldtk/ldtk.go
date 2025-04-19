package ldtk

import (
	"embed"
	"log"

	"github.com/TheBitDrifter/bappa/blueprint/ldtk"
)

//go:embed data.ldtk
var data embed.FS

var DATA = func() *ldtk.LDtkProject {
	project, err := ldtk.Parse(data, "../shared/ldtk/data.ldtk")
	if err != nil {
		log.Fatal(err)
	}
	return project
}()
