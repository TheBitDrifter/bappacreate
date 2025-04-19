
# Platformer Split Screen LDTK (Networked)Template

This is the Bappa Split Screen LDTK platformer template (Networked). In here you will find commented starter code for your project.

## Whats Included

By default this template provides one scene, up to 100 players, music, walking sounds, standard cameras that follows the players, basic player movement,
basic physics and collision resolution, one way platforms, split screen, multi scene support, and slope support.

This Project leverages LDTK to build its levels. The `data.ldtk` is located at `shared/ldtk/data.ldtk`

<https://ldtk.io/>

## Controls

- Movement: WASD (player one), arrow keys (player two)
- Toggle debug view: 0 key

## Asset Credits

- <https://tommusic.itch.io/free-fantasy-200-sfx-pack>
- <https://rustedstudio.itch.io/free-music-ambient-lofi-jazz-mp3-midi>

## Run Project

Due to Bappa's decoupled architecture the game/sim can be run in both networked and non networked 'modes'.

To run networked simply run the server and client:

```bash
cd yourproject/server
go mod tidy
go run .
```

```bash
cd yourproject/client
go mod tidy
go run .
```

To run offline:

```bash
cd yourproject/standalone
go mod tidy
go run .
```
