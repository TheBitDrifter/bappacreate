package ldtk

import (
	"embed"
	"log"

	blueprintldtk "github.com/TheBitDrifter/blueprint/ldtk"
)

//go:embed data.ldtk
var data embed.FS

var DATA = func() *blueprintldtk.LDtkProject {
	project, err := blueprintldtk.Parse(data, "./ldtk/data.ldtk")
	if err != nil {
		log.Fatal(err)
	}
	return project
}()
