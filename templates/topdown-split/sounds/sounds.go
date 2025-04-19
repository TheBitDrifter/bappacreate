package sounds

import "github.com/TheBitDrifter/bappa/blueprint/client"

var Run = client.SoundConfig{
	Path: "sounds/run.wav",
	// We could increase the count for two players, but it sounds too clustered in my opinion
	AudioPlayerCount: 1,
}

var Music = client.SoundConfig{
	Path:             "sounds/fantasy_music.wav",
	AudioPlayerCount: 1,
}
