package sounds

import "github.com/TheBitDrifter/bappa/blueprint/client"

var Run = client.SoundConfig{
	Path:             "sounds/run.wav",
	AudioPlayerCount: 1, // Too much running sounds clustered, keep at one
}

var Jump = client.SoundConfig{
	Path:             "sounds/jump.wav",
	AudioPlayerCount: 2, // matches max player count
}

var Land = client.SoundConfig{
	Path:             "sounds/land.wav",
	AudioPlayerCount: 2, // matches max player count
}

var Music = client.SoundConfig{
	Path:             "sounds/music.wav",
	AudioPlayerCount: 1,
}
