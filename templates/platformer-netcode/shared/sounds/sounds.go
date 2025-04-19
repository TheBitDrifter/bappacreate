package sounds

import "github.com/TheBitDrifter/bappa/blueprint/client"

var Run = client.SoundConfig{
	Path:             "sounds/run.wav",
	AudioPlayerCount: 1,
}

var Jump = client.SoundConfig{
	Path:             "sounds/jump.wav",
	AudioPlayerCount: 1,
}

var Land = client.SoundConfig{
	Path:             "sounds/land.wav",
	AudioPlayerCount: 1,
}

var Music = client.SoundConfig{
	Path:             "sounds/music.wav",
	AudioPlayerCount: 1,
}
