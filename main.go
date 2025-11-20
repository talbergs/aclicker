package main

import (
	"clicker2/game"
	"clicker2/shaders"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Game implements ebiten.Game interface.
type Game struct {
	state          *game.Game
	rockImage      *ebiten.Image
	offscreen1     *ebiten.Image // For shader chaining
	offscreen2     *ebiten.Image // For shader chaining
	rockPos        image.Point
	upgradeButton  image.Rectangle
	shadersEnabled bool
	time           float32
}

// Update proceeds the game state.
// Update is called every tick (1/60 second).
func (g *Game) Update() error {
	// Handle input
	g.handleInput()

	// Increment time for warp shader
	g.time++

	return nil
}

func (g *Game) handleInput() {
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
		if cursorPoint.In(g.upgradeButton) {
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
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the rock to an offscreen image
	g.offscreen1.Clear()
	op := &ebiten.DrawImageOptions{}
	g.offscreen1.DrawImage(g.rockImage, op)

	var finalImage *ebiten.Image
	if g.shadersEnabled {
		// 1. Grayscale
		g.offscreen2.Clear()
		g.offscreen2.DrawRectShader(g.rockImage.Bounds().Dx(), g.rockImage.Bounds().Dy(), shaders.GrayscaleShader, &ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{g.offscreen1},
		})

		// 2. Invert
		g.offscreen1.Clear()
		g.offscreen1.DrawRectShader(g.rockImage.Bounds().Dx(), g.rockImage.Bounds().Dy(), shaders.InvertShader, &ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{g.offscreen2},
		})

		// 3. Warp
		g.offscreen2.Clear()
		g.offscreen2.DrawRectShader(g.rockImage.Bounds().Dx(), g.rockImage.Bounds().Dy(), shaders.WarpShader, &ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{g.offscreen1},
			Uniforms: map[string]interface{}{
				"Time": g.time / 60.0, // Convert ticks to seconds
			},
		})
		finalImage = g.offscreen2
	} else {
		finalImage = g.offscreen1
	}

	// Draw the final rock image to the screen
	finalOp := &ebiten.DrawImageOptions{}
	finalOp.GeoM.Translate(float64(g.rockPos.X), float64(g.rockPos.Y))
	screen.DrawImage(finalImage, finalOp)


	// Draw the health bar
	g.drawHealthBar(screen)

	// Draw the upgrade button
	ebitenutil.DrawRect(screen, float64(g.upgradeButton.Min.X), float64(g.upgradeButton.Min.Y), float64(g.upgradeButton.Dx()), float64(g.upgradeButton.Dy()), color.RGBA{R: 100, G: 100, B: 100, A: 255})

	// Draw the stats
	msg := fmt.Sprintf("Rock Health: %d\nDust: %d\nDamage: %d\nUpgrade Cost: %d\nShaders: %t (Space)", g.state.TheRock.Health, g.state.ThePlayer.Dust, g.state.ThePlayer.Damage, game.UpgradeCost, g.shadersEnabled)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) drawHealthBar(screen *ebiten.Image) {
	barWidth := 100.0
	barHeight := 10.0
	barX := float64(g.rockPos.X) - (barWidth-float64(g.rockImage.Bounds().Dx()))/2
	barY := float64(g.rockPos.Y) - barHeight - 5 // 5 pixels above the rock

	healthPercentage := float64(g.state.TheRock.Health) / float64(game.InitialRockHealth)
	if healthPercentage < 0 {
		healthPercentage = 0
	}

	// Draw health bar background
	ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Draw health bar foreground
	ebitenutil.DrawRect(screen, barX, barY, barWidth*healthPercentage, barHeight, color.RGBA{R: 0, G: 255, B: 0, A: 255})
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return the screen size.
// For more detailed explanations, see https://github.com/hajimehoshi/ebiten/v2/wiki/Ebiten's-viewports.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
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

	// Define upgrade button
	buttonWidth := 150
	buttonHeight := 30
	upgradeButton := image.Rect(10, screenHeight-buttonHeight-10, 10+buttonWidth, screenHeight-10)

	game := &Game{
		state:          game.NewGame(),
		rockImage:      rockImage,
		offscreen1:     ebiten.NewImage(rockImage.Size()),
		offscreen2:     ebiten.NewImage(rockImage.Size()),
		rockPos:        image.Point{X: rockX, Y: rockY},
		upgradeButton:  upgradeButton,
		shadersEnabled: false,
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Clicker2")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
