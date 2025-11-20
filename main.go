package main

import (
	"clicker2/game"
	"clicker2/game/hud"
	"clicker2/shaders"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// EbitenGame implements ebiten.Game interface.
type EbitenGame struct {
	state          *game.Game
	hud            *hud.HUD
	rockImage      *ebiten.Image
	rockPos        image.Point
	shadersEnabled bool
	time           float32
}

// Update proceeds the game state.
// Update is called every tick (1/60 second).
func (g *EbitenGame) Update() error {
	// Handle input
	g.handleInput()

	// Quit game
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		log.Println("Game quit by user.")
		return ebiten.Termination
	}

	// Increment time for warp shader
	g.time++

	return nil
}

func (g *EbitenGame) handleInput() {
	// Mouse clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cursorPoint := image.Point{X: x, Y: y}

		// Check for rock click
		rockBounds := image.Rectangle{Min: g.rockPos, Max: g.rockPos.Add(g.rockImage.Bounds().Size())}
		if cursorPoint.In(rockBounds) {
			g.state.Click()
			return
		}

		// Check for upgrade click
		if cursorPoint.In(g.hud.UpgradeButton) {
			g.state.UpgradeDamage()
			return
		}
	}

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
	var finalImage *ebiten.Image
	if g.shadersEnabled {
		finalImage = shaders.Apply(g.rockImage,
			shaders.Grayscale(),
			shaders.Invert(),
			shaders.Warp(g.time/60.0),
		)
	} else {
		finalImage = g.rockImage
	}

	// Draw the final rock image to the screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.rockPos.X), float64(g.rockPos.Y))
	screen.DrawImage(finalImage, op)

	// Draw the health bar
	g.hud.DrawHealthBar(screen, g.rockPos, g.rockImage, g.state.TheRock.Health, game.InitialRockHealth)

	// Draw the HUD
	g.hud.Draw(screen, g.state.TheRock.Health, g.state.ThePlayer.Dust, g.state.ThePlayer.Damage, game.UpgradeCost, g.shadersEnabled)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return the screen size.
// For more detailed explanations, see https://github.com/hajimehoshi/ebiten/v2/wiki/Ebiten's-viewports.
func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Create a placeholder rock image
	rockImage := ebiten.NewImage(32, 32)
	rockImage.Fill(color.RGBA{R: 139, G: 69, B: 19, A: 255}) // SaddleBrown

	// Center the rock
	rockW, rockH := rockImage.Size()
	rockX := screenWidth/2 - rockW/2
	rockY := screenHeight/2 - rockH/2

	// Initialize HUD
	gameHUD := hud.NewHUD(screenWidth, screenHeight)

	game := &EbitenGame{
		state:          game.NewGame(),
		hud:            gameHUD,
		rockImage:      rockImage,
		rockPos:        image.Point{X: rockX, Y: rockY},
		shadersEnabled: false,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Clicker2")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
