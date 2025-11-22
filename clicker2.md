Dark genre game dynamics, a clicker.
You click the friendly rock and thus mine it's health.
The result of efficient clicking allows to upgrade stats making minig more efficient.

Game story:
- main state is 10000000 units of health shown in healthbar
- a click on rock given the accumulated compute gives approriate amount of damge and dust
- use the dust to buy various attack tools and skill coeificients improving mining process efficiency

Game State:
- `Game` struct: Holds the overall game state, including the rock, player, upgrades, event dispatcher, and various game flags (e.g., `IsPaused`, `GameOver`).
- `Rock` struct: Represents the clickable entity with its `Health`.
- `Player` struct: Represents the user's state with `Dust` and `Damage`.

Game inputs:
- q - quit
- s - save
- l - load
- h - toggle pause/shortcuts overview
- mouse-down - click event
- mouse-position - x,y (for click analysis and shader effects)
- scroll - adjust audio volume

Developer Inputs:
- space - toggle shaders
- F1 - set game state to early game
- F2 - set game state to mid game
- F3 - set game state to end game ready
- F4 - toggle debug auto-clicker
- F5 - cycle debug auto-clicker speed

Game view:
- 800x600px
- scroll on screen will adjust audio volume with visual feedback
- Rock sprite changes dynamically based on its health (full, cracked, shattered).
- HUD displays health bar, dust, and upgrade options.

Game technology:
- using golang and ebiten/v2 for 2d graphics
- rules of the game are structured in separate package that will encapsulate boundaries
- GPU is used for running various shaders, app decides them dynamically

Game assets:
- sprites.png file used for scene creation
-- capture hardcoded regions from sprite as variables "stone", "marketplace", "clouds" etc..
- game assets package also provides music player api and can be used in shaders

Game design - events driven (implemented using game/events and game/eventstore packages).
- actions (state transitions)
-- user stats update (i.e. increment damage by 10)
- events (facts)
-- user clicked stone
-- user stats update (an action example, were event is sourced from 2 past events)
--- user clicked stone for 10 times in past second
- audio stream is also passed inbetween shaders

Developer experience:
- fast tests fast feedback loop
- game tests are fully without any UI
- ./run.sh and ./build.sh scripts given
- `go mod tidy` is used to manage Go module dependencies.
- The `log` package is used for in-game logging and debugging.
