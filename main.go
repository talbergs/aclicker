package main

import (
	"clicker2/assets"
	"clicker2/game"
	"clicker2/game/clickanalysis"
	"clicker2/game/hud"
	"clicker2/shaders"
	"image"
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
	state             *game.Game
	hud               *hud.HUD
	clickGrid         *clickanalysis.ClickGrid
	rockImage         *ebiten.Image
	rockPos           image.Point
	marketplaceImage  *ebiten.Image
	marketplacePos    image.Point
	shadersEnabled    bool
	time              float32
	lastMouseX        float32
	lastMouseY        float32
	clickSpeed        float32
	lastClickPos      image.Point
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

	// Update click grid heat decay
	g.clickGrid.Update()

	g.clickSpeed *= 0.95

	return nil
}

func (g *EbitenGame) handleInput() {
	// Mouse clicks
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		cursorPoint := image.Point{X: x, Y: y}
		g.lastClickPos = cursorPoint
		g.clickSpeed += 0.5

		// Check for rock click
		rockBounds := image.Rectangle{Min: g.rockPos, Max: g.rockPos.Add(g.rockImage.Bounds().Size())}
		if cursorPoint.In(rockBounds) {
			g.state.Click()
			g.clickGrid.AddClick(cursorPoint.X, cursorPoint.Y, screenWidth, screenHeight)
			g.lastMouseX = float32(cursorPoint.X) / float32(screenWidth)
			g.lastMouseY = float32(cursorPoint.Y) / float32(screenHeight)
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
	// Draw the desert background
	if g.shadersEnabled {
		x, y := ebiten.CursorPosition()
		op := &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]interface{}{
				"Time":         g.time / 60.0,
				"Resolution":   []float32{screenWidth, screenHeight},
				"Mouse":        []float32{float32(x), float32(y)},
				"ClickSpeed":   g.clickSpeed,
				"LastClickPos": []float32{float32(g.lastClickPos.X), float32(g.lastClickPos.Y)},
			},
		}
		screen.DrawRectShader(screenWidth, screenHeight, shaders.DesertShader, op)
	}

	// Draw the marketplace
	opMarketplace := &ebiten.DrawImageOptions{}
	opMarketplace.GeoM.Translate(float64(g.marketplacePos.X), float64(g.marketplacePos.Y))
	screen.DrawImage(g.marketplaceImage, opMarketplace)

	var finalImage *ebiten.Image
	if g.shadersEnabled {
		clickGridEbitenImage := ebiten.NewImageFromImage(g.clickGrid.ToRGBA())
		finalImage = shaders.Apply(g.rockImage, clickGridEbitenImage,
			// shaders.Grayscale(),
			// shaders.Invert(),
			// shaders.Warp(g.time/60.0),
			shaders.TimeClick(g.time/60.0),
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
	// Center the rock
	rockW, rockH := assets.RockSprite.Size()
	rockX := screenWidth/2 - rockW/2
	rockY := screenHeight/2 - rockH/2

	// Position the marketplace
	marketplaceX := 50
	marketplaceY := screenHeight/2 - 128/2

	// Initialize HUD
	gameHUD := hud.NewHUD(screenWidth, screenHeight)

	// Initialize ClickGrid
	clickGrid := clickanalysis.NewClickGrid(rockW, rockH)

	game := &EbitenGame{
		state:             game.NewGame(),
		hud:               gameHUD,
		clickGrid:         clickGrid,
		rockImage:         assets.RockSprite,
		rockPos:           image.Point{X: rockX, Y: rockY},
		marketplaceImage:  assets.MarketplaceSprite,
		marketplacePos:    image.Point{X: marketplaceX, Y: marketplaceY},
		shadersEnabled:    true,
		clickSpeed:        0.0,
		lastClickPos:      image.Point{X: 0, Y: 0},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Clicker2")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
