# Bappa Game Template Generator

A simple tool for generating new Bappa game projects from templates.
<img width="800" alt="lol" src="https://github.com/user-attachments/assets/1cb4159b-5761-4635-b0f3-040b27264736" />

## Overview

Bappa Game Template Generator helps you quickly bootstrap new game projects using the Bappa game engine. It creates a complete project structure with all necessary files and dependencies based on predefined templates.

## Installation

### Option 1: Install with Go (Recommended)

```bash
go install github.com/TheBitDrifter/bappacreate@latest
```

### Option 2: Install from source

```bash
# Clone the repository
git clone https://github.com/TheBitDrifter/bappacreate.git
cd bappacreate
# Build and install the binary
go install
```

### After Installation

After installing with `go install`, make sure the Go bin directory is in your PATH:

1. Find your Go binary path:

   ```bash
   go env GOPATH
   ```

2. Add the Go bin directory to your PATH:
   For bash (add to ~/.bashrc):

   ```bash
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

   For zsh (add to ~/.zshrc):

   ```bash
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

   Then reload your shell configuration:

   ```bash
   source ~/.bashrc   # For bash
   source ~/.zshrc    # For zsh
   ```

3. Verify the installation:

   ```bash
   bappacreate --help
   ```

## Usage

```bash
bappacreate username/project-name [--template <template-name>]
```

The `username/` prefix is important as it will be used to create the proper Go module path (`github.com/username/project-name`).

### Examples

Create a top-down game (default template):

```bash
bappacreate johndoe/my-awesome-game
```

Create a platformer game:

```bash
bappacreate johndoe/my-platformer --template platformer
```

Create a platformer game with LDtk support:

```bash
bappacreate johndoe/my-ldtk-platformer --template platformer-ldtk
```

Create a sandbox game:

```bash
bappacreate johndoe/my-sandbox-world --template sandbox
```

Create a networked platformer game:

```bash
bappacreate johndoe/my-netcode-game --template platformer-netcode
```

## Available Templates

| Template | Description |
|----------|-------------|
| `topdown` | A top-down perspective game (default) |
| `topdown-split` | A top-down game with split-screen co-op support |
| `platformer` | A simple platformer game |
| `platformer-split` | A platformer game with split-screen co-op support |
| `platformer-ldtk` | A platformer game with LDtk level editor support |
| `platformer-split-ldtk` | A platformer game with split-screen co-op and LDtk support |
| `platformer-netcode` | A platformer game with multiplayer netcode support |
| `sandbox` | An open sandbox game environment |

### Special Note for Netcode Template

For the `platformer-netcode` template, you need to access the template directory. When installing via `go install`, you'll need to clone the repository:

```bash
git clone https://github.com/TheBitDrifter/bappacreate.git
cd bappacreate
go build .
./bappacreate username/project-name --template platformer-netcode
```

You can also skip using bappacreate all together and simple grab the files at `/templates/platformer-netcode`.

Since the server/client are distinct modules(with their own deps and .mod file) they cannot be embedded as part of the
bappacreate binary, making it a nuisance to implement via this installer.

### Template Structure

- **Standard Templates** (`topdown`, `platformer`, `sandbox`): Single-player games with a standard view.
- **Split Templates** (`topdown-split`, `platformer-split`, `platformer-split-ldtk`): Games with split-screen co-op multiplayer support, allowing two or more players to play simultaneously on the same screen.
- **LDtk Templates** (`platformer-ldtk`, `platformer-split-ldtk`): Games that support the LDtk level editor for easier level design.
- **Netcode Template** (`platformer-netcode`): A client/server architecture for multiplayer networked games.

## Project Structure

When you create a new project, the following structure will be generated:

```
my-awesome-game/
├── assets/
│   ├── images/
│   └── sounds/
├── players/        # Only for split-screen co-op templates
├── splitscreen/    # Only for split-screen co-op templates
├── go.mod
├── go.sum
└── main.go
```

For the netcode template, a different structure is created:

```
my-netcode-game/
├── client/         # Client-side code
├── server/         # Server-side code
├── shared/         # Code shared between client and server
├── sharedclient/   # Code shared between client and standalone
│   └── assets/
│       ├── images/
│       └── sounds/
├── bot/            # AI bot implementation
└── standalone/     # Single-player version
```

## After Creating a Project

After creating your project, navigate to the project directory and run it:

```bash
cd my-awesome-game
go run .
```

For netcode projects:

```bash
# Run server
cd my-netcode-game/server
go mod tidy
go run .

# Run client in another terminal
cd my-netcode-game/client
go mod tidy
go run .

# Or run standalone version
cd my-netcode-game/standalone
go mod tidy
go run .
```

## Contributing

We welcome contributions to add new templates or improve existing ones! Simply add your template to the `templates/` directory and submit a pull request.

## License

[MIT License](LICENSE)
