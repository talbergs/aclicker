# CLAUDE.md - AI Assistant Guide for Clicker2 Project

## Project Overview

**Clicker2** is a dark genre clicker game built in Go using the Ebiten 2D game engine. The game features a unique narrative where players click a "friendly rock" to mine dust, but each click depletes the rock's health. The ultimate goal is to reach a point where the player must decide between consuming the rock entirely or choosing to let it rest - creating a moral dilemma uncommon in clicker games.

### Core Game Mechanics
- Click the friendly rock to mine "dust" (currency)
- Each click damages the rock's health
- Spend dust on upgrades to mine more efficiently
- Progress through upgrade tiers that reveal the rock's sentience
- Reach an endgame choice: "Take the Heart" (bad ending) or "Let it Rest" (good ending)

### Key Technologies
- **Language**: Go 1.24.3
- **Game Engine**: Ebiten v2.9.4 (2D game library)
- **Graphics**: GPU shaders (Kage language), sprite-based rendering
- **Audio**: Built-in Ebiten audio with MP3 support
- **Architecture**: Event-driven with event sourcing

## Codebase Structure

```
aclicker/
├── main.go                      # Entry point, game loop, rendering
├── game/                        # Core game logic package
│   ├── game.go                  # Main game state, Rock, Player structs
│   ├── game_test.go             # Comprehensive game logic tests
│   ├── upgrades.go              # Upgrade system and manager
│   ├── events/                  # Event-driven architecture
│   │   └── events.go            # Event interfaces, dispatcher, handlers
│   ├── eventstore/              # Event sourcing implementation
│   │   └── eventstore.go        # File-based event store (events.log)
│   ├── hud/                     # Heads-Up Display
│   │   └── hud.go               # Stats display, health bar, upgrade buttons
│   ├── clickanalysis/           # Click heat tracking
│   │   └── clickanalysis.go    # Grid-based click visualization
│   └── errors/                  # Custom error types
│       └── errors.go            # Game-specific errors
├── shaders/                     # GPU shader effects
│   ├── shaders.go               # Shader loading and compilation
│   ├── chain.go                 # Shader effect chaining utility
│   └── kage/                    # Kage shader source files
│       ├── desert.kage          # Procedural nebula background
│       ├── clickheat.kage       # Click heat visualization
│       ├── grayscale.kage       # Grayscale effect
│       ├── invert.kage          # Color inversion
│       ├── timeclick.kage       # Time-based click effect
│       └── warp.kage            # Warp distortion effect
├── assets/                      # Game assets package
│   ├── assets.go                # Asset loading and management
│   ├── sprite.png               # Main sprite sheet
│   └── *.mp3                    # Audio files (music, SFX)
├── *.md                         # Documentation files
│   ├── GEMINI.md                # Gemini AI assistant guide
│   ├── architecture.md          # System architecture diagram
│   ├── codebase_analysis.md     # Detailed codebase analysis
│   ├── game_mechanics.md        # Game mechanics and design
│   ├── testing_instructions.md  # Comprehensive testing guide
│   ├── agentic_development.md   # Development process guide
│   └── iterations.md            # Development iteration log
├── build.sh                     # Build script
├── run.sh                       # Run script (removes save first)
└── test.sh                      # Test runner script
```

## Architecture and Design Patterns

### 1. Event-Driven Architecture with Event Sourcing

**Critical Pattern**: The game uses a strict event-driven architecture where all state changes are represented as events.

- **Events**: Factual occurrences (e.g., "user clicked stone")
- **Actions/Handlers**: Functions that apply events to modify state
- **Event Store**: All events are persisted to `events.log` for replay

**Key Benefits**:
- Predictable, testable game logic
- Robust save/load through event replay
- Clear separation between actions and state changes

**AI Assistant Guidelines**:
- ✅ **ALWAYS** dispatch events for state changes, never modify state directly
- ✅ Create new event types in `game/events/events.go` when adding features
- ✅ Register event handlers in `NewGame()` function
- ❌ **NEVER** directly modify `TheRock.Health`, `ThePlayer.Dust`, etc. - use events

Example:
```go
// ❌ WRONG - Direct state modification
g.ThePlayer.Dust += 10

// ✅ CORRECT - Event-driven approach
g.Dispatcher.Dispatch(&events.DustEarnedEvent{Amount: 10})
```

### 2. Separation of Concerns

The codebase maintains clear boundaries:
- `game/` - Pure game logic, no rendering code
- `main.go` - Presentation layer, rendering, input handling
- `shaders/` - Visual effects, independent of game logic
- `assets/` - Resource management

**AI Assistant Guidelines**:
- ✅ Keep game logic in `game/` package
- ✅ Keep rendering/drawing code in `main.go`
- ❌ Don't add Ebiten-specific code to `game/` package

### 3. Procedural Content Generation

The game uses GPU shaders for dynamic backgrounds:
- `desert.kage` generates animated nebula effects
- Shaders respond to game state (health percentage, click speed, mouse position)

### 4. Save/Load System

Two mechanisms:
1. **JSON Save**: Traditional save file (`save.json`) with game state snapshot
2. **Event Sourcing**: Can reconstruct game state by replaying `events.log`

## Development Workflow

### Building and Running

```bash
# Build the game
./build.sh
# Or manually:
go build -o clicker2 .

# Run the game (removes save.json first)
./run.sh
# Or manually:
rm -f save.json && ./clicker2

# Run tests
./test.sh
# Or manually:
go test ./...
```

### Key Commands and Controls

**In-Game Controls**:
- `Mouse Click` - Click the rock
- `Mouse Scroll` - Adjust volume (up/down)
- `S` - Save game state
- `L` - Load game state
- `Q` - Quit game
- `P` - Pause game (if implemented)

### Testing Approach

**Location**: `game/game_test.go`

The project uses comprehensive table-driven tests with:
- Event replay testing
- Boundary condition validation
- State assertion checks
- Upgrade mechanics verification

**AI Assistant Guidelines**:
- ✅ Add tests for new features in `game/game_test.go`
- ✅ Use table-driven test patterns
- ✅ Test boundary conditions (e.g., health at 0, exact upgrade cost)
- ✅ Verify event handlers correctly apply state changes

See `testing_instructions.md` for detailed testing procedures.

## Key Conventions for AI Assistants

### Code Style

1. **Go Standard**: Follow standard Go conventions (gofmt, proper error handling)
2. **Package Organization**:
   - Game logic in `game/` package
   - Each subpackage serves a single purpose
3. **Constants**: Define constants at package level (e.g., `InitialRockHealth`)
4. **Exported vs Unexported**: Use proper capitalization for visibility

### Adding New Features

When adding new game features, follow this checklist:

1. **Define Event Type** (in `game/events/events.go`):
   ```go
   type NewFeatureEvent struct {
       // Event data
   }

   func (e *NewFeatureEvent) EventType() string {
       return "NewFeature"
   }
   ```

2. **Create Event Handler** (in appropriate `game/*.go` file):
   ```go
   func (g *Game) ApplyNewFeatureEvent(event events.Event) {
       e := event.(*events.NewFeatureEvent)
       // Apply state changes here
   }
   ```

3. **Register Handler** (in `game/game.go` `NewGame()` function):
   ```go
   g.Dispatcher.Register("NewFeature", g.ApplyNewFeatureEvent)
   ```

4. **Dispatch Event** (where action occurs):
   ```go
   g.Dispatcher.Dispatch(&events.NewFeatureEvent{...})
   ```

5. **Add Tests** (in `game/game_test.go`)

### Upgrade System

Upgrades are defined in `game/upgrades.go`:
- Each upgrade has ID, name, base cost, cost multiplier, max level, and effect function
- Use `UpgradeManager` to manage upgrade state
- Apply upgrade effects through events

**Current Upgrades**:
- `stronger_pickaxe` - Increases damage per click
- `auto_clicker_v0_1` - Enables auto-clicking
- `auto_clicker_v1_0` - Faster auto-clicking (permanent)
- `heart_of_the_mountain` - Triggers endgame choice

### Visual and Audio Guidelines

**Sprites**:
- All sprites in single `assets/sprite.png` file
- Sprite coordinates defined in `assets/assets.go`
- Load sprites using `assets` package

**Shaders**:
- Written in Kage language (`.kage` files)
- Embedded at compile time using `//go:embed`
- Compiled in `shaders/shaders.go`
- Pass game state to shaders via uniforms (Time, Mouse, ClickSpeed, HealthPercentage)

**Audio**:
- Music crossfades based on rock health (healthy ↔ melancholic)
- SFX for clicks, upgrades, errors
- Volume control via mouse scroll

## Important Files to Reference

### Before Making Changes

Always reference these files to understand context:

1. **`clicker2.md`** - Project roadmap (referenced by GEMINI.md)
2. **`game_mechanics.md`** - Complete game design document
3. **`architecture.md`** - System architecture diagram
4. **`codebase_analysis.md`** - Design coherence analysis
5. **`testing_instructions.md`** - Testing procedures and assertions

### Core Implementation Files

1. **`main.go`** - Game loop, rendering, input handling
2. **`game/game.go`** - Core game state and logic
3. **`game/events/events.go`** - Event system
4. **`game/upgrades.go`** - Upgrade definitions
5. **`game/hud/hud.go`** - UI rendering

## Common Tasks and Patterns

### Adding a New Upgrade

1. Add upgrade definition in `game/upgrades.go` `Init()` method
2. Create event type if upgrade has unique effects
3. Add event handler for upgrade purchase
4. Update HUD to display new upgrade button
5. Add test cases

### Modifying Game State

```go
// Define what changed (event)
type StateChangeEvent struct {
    NewValue int
}

func (e *StateChangeEvent) EventType() string {
    return "StateChange"
}

// Define how to apply it (handler)
func (g *Game) ApplyStateChangeEvent(event events.Event) {
    e := event.(*StateChangeEvent)
    g.SomeField = e.NewValue
}

// Register in NewGame()
g.Dispatcher.Register("StateChange", g.ApplyStateChangeEvent)

// Use it
g.Dispatcher.Dispatch(&StateChangeEvent{NewValue: 42})
```

### Adding Visual Effects

1. Create `.kage` shader file in `shaders/kage/`
2. Embed shader in `shaders/shaders.go`
3. Compile shader at initialization
4. Pass uniforms from `main.go` to control effect
5. Apply in `Draw()` method

## Error Handling

- Custom errors in `game/errors/errors.go`
- Check `Dispatcher.Dispatch()` returns for event errors
- Validate upgrade purchases return errors for insufficient funds
- Handle file I/O errors for save/load

## Performance Considerations

- Shaders run on GPU for efficient rendering
- Event store appends to file (potential performance issue for long sessions)
- Click grid uses decay mechanism to manage memory
- Sprite sheet reduces asset loading overhead

## Development History

See `iterations.md` for chronological development log.

Key milestones:
- Nov 20, 2025: Core game state, click mechanics, health bar
- Nov 20, 2025: Upgrades, save/load system
- Recent: Event-driven architecture, event sourcing, comprehensive testing

## AI Assistant Best Practices

### DO:
- ✅ Read referenced documentation files before making changes
- ✅ Use event-driven patterns for all state changes
- ✅ Write comprehensive tests with boundary conditions
- ✅ Follow Go conventions and existing code style
- ✅ Keep game logic separate from presentation
- ✅ Add comments for complex shader or game logic
- ✅ Test changes by running `./test.sh`
- ✅ Verify builds with `./build.sh`

### DON'T:
- ❌ Modify game state directly without events
- ❌ Add Ebiten dependencies to `game/` package
- ❌ Create new files without clear purpose
- ❌ Ignore existing upgrade/event patterns
- ❌ Skip tests when adding features
- ❌ Hardcode values that should be constants
- ❌ Mix rendering code with game logic

## Debugging Tips

1. **Event Log**: Check `events.log` to see event history
2. **Save File**: Inspect `save.json` for current game state
3. **Tests**: Run specific tests: `go test ./game -run TestName`
4. **Build Errors**: Shader compilation errors appear at runtime
5. **Game State**: Add debug prints in event handlers to trace state changes

## Contact and Resources

- **Go Documentation**: https://golang.org/doc/
- **Ebiten Documentation**: https://ebitengine.org/en/documents/
- **Kage Shader Guide**: https://ebitengine.org/en/documents/shader.html

---

**Last Updated**: 2025-12-03
**Project Version**: Active Development
**Maintained By**: AI-assisted development workflow
