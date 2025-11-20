Always refer to the project_roadmap_file (./clicker2.md) and sections in it before considering further the decisions you will make upon my prompts. Keep this file as project roadmap.

# Project: Clicker2

## Project Overview

This project is a dark genre clicker game, as outlined in `clicker2.md`. The core mechanic involves clicking a "friendly rock" to mine its health, gather "dust", and use it to purchase upgrades for more efficient mining.

*   **Main Technologies:** Go, a 2D graphics library (unspecified), and GLSL shaders.
*   **Architecture:** The game uses an event-driven architecture. The core game logic and rules are to be encapsulated in a separate Go package to maintain clear boundaries.

## Building and Running

The project roadmap in `clicker2.md` does not specify the build and run commands. Assuming a standard Go project structure, the following commands are likely to be used.

**TODO:** Verify and update these commands once the project structure is in place.

```bash
# Build the project
go build -o clicker2 .

# Run the game
./clicker2

# Run tests (once tests are added)
go test ./...
```

## Development Conventions

Based on `clicker2.md`, the following conventions should be followed:

*   **Event-Driven Design:** The game state is modified through actions (state transitions) which are triggered by events.
    *   **Events:** Factual occurrences, e.g., `user clicked stone`.
    *   **Actions:** State changes that result from events, e.g., `user stats update`.
*   **Game Logic Encapsulation:** All game rules, boundaries, and core logic should be placed within a dedicated Go package.
*   **Asset Management:**
    *   A single `sprites.png` file will be used for all game sprites. Specific regions of this sprite sheet will be mapped to variables like "stone", "marketplace", etc.
    *   A dedicated "game assets" package will provide an API for music and sound.
*   **View:** The game window will be 800x600px.
*   **Inputs:**
    *   `q`: quit
    *   `s`: save
    *   `l`: load
    *   `mouse-down`: click event
    *   `scroll`: adjust audio volume