package main

import (
	"clicker2/assets"
	"clicker2/game"
	"clicker2/game/clickanalysis"
	"clicker2/game/hud"
	"clicker2/shaders"
	"image"
	"image/color" // Import for color
	"log"

	"golang.org/x/image/font/basicfont" // Import for basic font
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text" // Import for text drawing
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Debug holds the game state for debugging purposes.
type Debug struct {
	state *game.Game
}

// NewDebug creates a new Debug instance.
func NewDebug(state *game.Game) *Debug {
	return &Debug{state: state}
}

// UpdateDebugAutoClicker updates the debug auto-clicker.
func (d *Debug) UpdateDebugAutoClicker() {
	// TODO: Implement auto-clicker logic
}

// HandleDeveloperKeybinds handles developer keybinds.
func (d *Debug) HandleDeveloperKeybinds() {
	// TODO: Implement developer keybinds
}

// EbitenGame implements ebiten.Game interface.
type EbitenGame struct {
	state             *game.Game
	hud               *hud.HUD
	clickGrid         *clickanalysis.ClickGrid
	currentRockSprite *ebiten.Image // Dynamically selected rock sprite
	rockPos           image.Point
	marketplaceImage  *ebiten.Image
	marketplacePos    image.Point
	shadersEnabled    bool
	time              float32
	lastMouseX        float32
	lastMouseY        float32
	clickSpeed        float32
	lastClickPos      image.Point
	debug             *Debug // New field for debug functionality
	IsPaused          bool   // New field to track if the game is paused
	ShowShortcuts     bool   // New field to track if shortcuts are displayed
}

// Update proceeds the game state.
// Update is called every tick (1/60 second).
func (g *EbitenGame) Update() error {
	// Handle input
	g.handleInput()

	// Quit game
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) || g.state.ShouldExit {
		log.Println("Game quit by user.")
		return ebiten.Termination
	}

	if g.IsPaused {
		return nil // Skip all game logic updates if paused
	}

	// Increment time for warp shader
	g.time++

	// Update click grid heat decay
	g.clickGrid.Update()

	// Update rock message timer
	if g.state.RockMessageTimer > 0 {
		g.state.RockMessageTimer -= 1.0 / float64(ebiten.TPS())
		if g.state.RockMessageTimer <= 0 {
			g.state.CurrentRockMessage = ""
		}
	}

	// Update music crossfade based on rock health
	healthPercentage := float64(g.state.TheRock.Health) / float64(game.InitialRockHealth)
	healthyVolume := healthPercentage
	melancholicVolume := 1.0 - healthPercentage

	// Clamp volumes to ensure they are between 0 and 1
	if healthyVolume < 0 { healthyVolume = 0 }
	if healthyVolume > 1 { healthyVolume = 1 }
	if melancholicVolume < 0 { melancholicVolume = 0 }
	if melancholicVolume > 1 { melancholicVolume = 1 }

	assets.HealthyMusicPlayer.SetVolume(healthyVolume)
	assets.MelancholicMusicPlayer.SetVolume(melancholicVolume)

	g.debug.UpdateDebugAutoClicker() // Call the debug package's auto-clicker update

	g.clickSpeed *= 0.95

	return nil
}

func (g *EbitenGame) handleInput() {
	// Toggle pause and shortcut display with KeyH
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		g.IsPaused = !g.IsPaused
		g.ShowShortcuts = !g.ShowShortcuts // Toggle shortcut display along with pause
		if g.IsPaused {
			log.Println("Game Paused. Press H for shortcuts.")
		} else {
			log.Println("Game Resumed.")
		}
	}

	// If the game is paused, only allow KeyH and KeyQ to be processed.
	// All other game-related inputs should be ignored.
	if g.IsPaused {
		return
	}

	g.handleMouseInput()
	g.handleGameKeybinds()
	g.debug.HandleDeveloperKeybinds() // Call the debug package's handler
}

func (g *EbitenGame) handleMouseInput() {
	// Mouse clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cursorPoint := image.Point{X: x, Y: y}
		g.lastClickPos = cursorPoint
		g.clickSpeed += 0.5

		// Check for rock click
		rockBounds := image.Rectangle{Min: g.rockPos, Max: g.rockPos.Add(g.currentRockSprite.Bounds().Size())}
		if cursorPoint.In(rockBounds) {
			g.state.Click()
			g.clickGrid.AddClick(cursorPoint.X, cursorPoint.Y, screenWidth, screenHeight)
			g.lastMouseX = float32(cursorPoint.X) / float32(screenWidth)
			g.lastMouseY = float32(cursorPoint.Y) / float32(screenHeight)
			assets.ClickSFXPlayer.Rewind()
			// assets.ClickSFXPlayer.Play()
			return
		}

		// Check for upgrade click
		if clickedUpgradeID := g.hud.GetClickedUpgradeID(cursorPoint); clickedUpgradeID != "" {
			if err := g.state.PurchaseUpgrade(clickedUpgradeID); err != nil {
				log.Printf("Error purchasing upgrade %s: %v", clickedUpgradeID, err)
				assets.ErrorSFXPlayer.Rewind()
				// assets.ErrorSFXPlayer.Play()
			} else {
				assets.UpgradeSFXPlayer.Rewind()
				// assets.UpgradeSFXPlayer.Play()
			}
			return
		}	}

	// Handle end-game choices
	if g.state.EndGameChoicePending {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			cursorPoint := image.Point{X: x, Y: y}
			if clickedChoiceID := g.hud.GetClickedChoiceID(cursorPoint); clickedChoiceID != "" {
				if clickedChoiceID == "take_heart" {
					g.state.TakeHeart()
				} else if clickedChoiceID == "let_rest" {
					g.state.LetRest()
				}
				return
			}
		}
	}

	// Adjust audio volume with scroll wheel
	_, wy := ebiten.Wheel()
	if wy != 0 {
		currentVolume := assets.HealthyMusicPlayer.Volume() // Assuming both players have same volume
		newVolume := currentVolume + wy*0.05 // Adjust sensitivity as needed

		if newVolume < 0 {
			newVolume = 0
		}
		if newVolume > 1 {
			newVolume = 1
		}

		assets.HealthyMusicPlayer.SetVolume(newVolume)
		assets.MelancholicMusicPlayer.SetVolume(newVolume)
	}
}

func (g *EbitenGame) handleGameKeybinds() {
	// Save game
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if err := g.state.Save(); err != nil {
			log.Printf("Error saving game: %v", err)
		} else {
			log.Println("Game saved!")
		}
	}

	// Load game
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		if err := g.state.Load(); err != nil {
			log.Printf("Error loading game: %v", err)
		} else {
			log.Println("Game loaded!")
		}
	}

	// Toggle shaders
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.shadersEnabled = !g.shadersEnabled
	}
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60 second).
func (g *EbitenGame) Draw(screen *ebiten.Image) {
	// Draw the desert background
	if g.shadersEnabled {
		x, y := ebiten.CursorPosition()
		healthPercentage := float32(g.state.TheRock.Health) / float32(game.InitialRockHealth)
		op := &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]interface{}{
				"Time":         g.time / 60.0,
				"Resolution":   []float32{screenWidth, screenHeight},
				"Mouse":        []float32{float32(x), float32(y)},
				"ClickSpeed":   g.clickSpeed,
				"LastClickPos": []float32{float32(g.lastClickPos.X), float32(g.lastClickPos.Y)},
				"HealthPercentage": healthPercentage,
			},
		}
		screen.DrawRectShader(screenWidth, screenHeight, shaders.DesertShader, op)
	}

	// Draw the marketplace
	opMarketplace := &ebiten.DrawImageOptions{}
	opMarketplace.GeoM.Translate(float64(g.marketplacePos.X), float64(g.marketplacePos.Y))
	screen.DrawImage(g.marketplaceImage, opMarketplace)

	var currentRockSprite *ebiten.Image
	healthPercentage := float64(g.state.TheRock.Health) / float64(game.InitialRockHealth)

	if healthPercentage > 0.75 {
		currentRockSprite = assets.RockSpriteFull
	} else if healthPercentage > 0.50 {
		currentRockSprite = assets.RockSpriteCracked1
	} else if healthPercentage > 0.25 {
		currentRockSprite = assets.RockSpriteCracked2
	} else {
		currentRockSprite = assets.RockSpriteShattered
	}
	g.currentRockSprite = currentRockSprite // Assign to struct field

	var finalImage *ebiten.Image
	if g.shadersEnabled {
		clickGridEbitenImage := ebiten.NewImageFromImage(g.clickGrid.ToRGBA())
		finalImage = shaders.Apply(currentRockSprite, clickGridEbitenImage,
			shaders.Grayscale(),
			shaders.Invert(),
			shaders.Warp(g.time/60.0),
			shaders.TimeClick(g.time/60.0),
		)
	} else {
		finalImage = currentRockSprite
	}

	// Draw the final rock image to the screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.rockPos.X), float64(g.rockPos.Y))
	screen.DrawImage(finalImage, op)

	// Draw the health bar
	g.hud.DrawHealthBar(screen, g.rockPos, currentRockSprite, g.state.TheRock.Health, game.InitialRockHealth)

	// Draw the HUD
	g.hud.Draw(screen, g.state, g.shadersEnabled)

	// Draw shortcuts if enabled
	if g.ShowShortcuts {
		// Draw a semi-transparent background
		overlay := ebiten.NewImage(screenWidth, screenHeight)
		overlay.Fill(color.RGBA{0, 0, 0, 180}) // Dark, semi-transparent
		screen.DrawImage(overlay, &ebiten.DrawImageOptions{})

		// Define font and color
		f := basicfont.Face7x13
		col := color.White
		lineHeight := f.Metrics().Height.Ceil() + 4 // Add some padding between lines

		// Starting position for text
		xOffset := 50
		yOffset := 50

		for i, shortcut := range g.Shortcuts() {
			text.Draw(screen, shortcut, f, xOffset, yOffset+(i*lineHeight), col)
		}
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return the screen size.
// For more detailed explanations, see https://github.com/hajimehoshi/ebiten/v2/wiki/Ebiten's-viewports.
func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Shortcuts returns a slice of strings, each representing a game shortcut.
func (g *EbitenGame) Shortcuts() []string {
	return []string{
		"--- Game Shortcuts ---",
		"H: Toggle Pause/Shortcuts",
		"Q: Quit Game",
		"S: Save Game",
		"L: Load Game",
		"Scroll: Adjust Music Volume",
		"Space: Toggle Shaders",
		"--- Developer Shortcuts ---",
		"F1: Set State Early Game",
		"F2: Set State Mid Game",
		"F3: Set State End Game Ready",
		"F4: Toggle Debug Auto-Clicker",
		"F5: Cycle Debug Auto-Clicker Speed",
	}
}

func main() {
	// Center the rock
	rockW, rockH := assets.RockSpriteFull.Size()
	rockX := screenWidth/2 - rockW/2
	rockY := screenHeight/2 - rockH/2

	// Position the marketplace
	marketplaceX := 50
	marketplaceY := screenHeight/2 - 128/2

	// Initialize game state
	gameState := game.NewGame()

	// Initialize HUD
	gameHUD := hud.NewHUD(screenWidth, screenHeight, gameState.Upgrades)

	// Initialize ClickGrid
	clickGrid := clickanalysis.NewClickGrid(rockW, rockH)

	// Initialize Debug
	debug := NewDebug(gameState)

	game := &EbitenGame{
		state:             gameState, // Use the already initialized gameState
		hud:               gameHUD,
		clickGrid:         clickGrid,
		rockPos:           image.Point{X: rockX, Y: rockY}, // Corrected rockY
		marketplaceImage:  assets.MarketplaceSprite,
		marketplacePos:    image.Point{X: marketplaceX, Y: marketplaceY},
		shadersEnabled:    true,
		clickSpeed:        0.0,
		lastClickPos:      image.Point{X: 0, Y: 0},
		debug:             debug, // Assign the new Debug instance
	}

	// Start background music
	assets.HealthyMusicPlayer.SetVolume(0.5) // Ensure initial volume
	// assets.HealthyMusicPlayer.Play()

	assets.MelancholicMusicPlayer.SetVolume(0.0) // Start muted

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Clicker2")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
