module github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/standalone

go 1.24.1

// replace github.com/TheBitDrifter/bappa/table => ../../../../Bappa/table/

// replace github.com/TheBitDrifter/bappa/warehouse => ../../../../Bappa/warehouse/

// replace github.com/TheBitDrifter/bappa/tteokbokki => ../../../../Bappa/tteokbokki/

// replace github.com/TheBitDrifter/bappa/blueprint => ../../../../Bappa/blueprint/

// replace github.com/TheBitDrifter/bappa/coldbrew => ../../../../Bappa/coldbrew/

// replace github.com/TheBitDrifter/bappa/environment => ../../../../Bappa/environment/

// replace github.com/TheBitDrifter/bappa/drip => ../../../../Bappa/drip/

replace github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared => ../shared/

replace github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/sharedclient => ../sharedclient/

require (
	github.com/TheBitDrifter/bappa/coldbrew v0.0.0-20250420132432-5606172c9a41
	github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/shared v0.0.0-00010101000000-000000000000
	github.com/TheBitDrifter/bappacreate/templates/platformer-netcode/sharedclient v0.0.0-00010101000000-000000000000
	github.com/hajimehoshi/ebiten/v2 v2.8.7
)

require (
	github.com/TheBitDrifter/bappa/blueprint v0.0.0-20250420132432-5606172c9a41 // indirect
	github.com/TheBitDrifter/bappa/drip v0.0.0-20250420132432-5606172c9a41 // indirect
	github.com/TheBitDrifter/bappa/environment v0.0.0-20250420132432-5606172c9a41 // indirect
	github.com/TheBitDrifter/bappa/table v0.0.0-20250420132432-5606172c9a41 // indirect
	github.com/TheBitDrifter/bappa/tteokbokki v0.0.0-20250420132432-5606172c9a41 // indirect
	github.com/TheBitDrifter/bappa/warehouse v0.0.0-20250420132432-5606172c9a41 // indirect
	github.com/TheBitDrifter/bark v0.0.0-20250302175939-26104a815ed9 // indirect
	github.com/TheBitDrifter/mask v0.0.1-early-alpha.1 // indirect
	github.com/TheBitDrifter/util v0.0.0-20241102212109-342f4c0a810e // indirect
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/oto/v3 v3.3.3 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/go-text/typesetting v0.2.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/image v0.25.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)
