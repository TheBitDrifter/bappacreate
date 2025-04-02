package sounds

import "github.com/TheBitDrifter/bappa/blueprint/client"

var Run = client.SoundConfig{
	Path:             "run.wav",
	AudioPlayerCount: 1, // Too much running sounds clustered, keep at one
}

var Jump = client.SoundConfig{
	Path:             "jump.wav",
	AudioPlayerCount: 2, // Increase to player count
}

var Land = client.SoundConfig{
	Path:             "land.wav",
	AudioPlayerCount: 2, // Increase to player count
}

var Music = client.SoundConfig{
	Path:             "music.wav",
	AudioPlayerCount: 1,
}
