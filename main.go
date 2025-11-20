package main

import (
	"clicker2/game"
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
	state     *game.Game
	rockImage *ebiten.Image
	rockPos   image.Point
}

// Update proceeds the game state.
// Update is called every tick (1/60 second).
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		rockBounds := image.Rectangle{Min: g.rockPos, Max: g.rockPos.Add(g.rockImage.Bounds().Size())}
		if (image.Point{X: x, Y: y}).In(rockBounds) {
			g.state.Click()
		}
	}
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60 second).
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the rock
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.rockPos.X), float64(g.rockPos.Y))
	screen.DrawImage(g.rockImage, op)

	// Draw the stats
	msg := fmt.Sprintf("Rock Health: %d\nDust: %d", g.state.TheRock.Health, g.state.ThePlayer.Dust)
	ebitenutil.DebugPrint(screen, msg)
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

	game := &Game{
		state:     game.NewGame(),
		rockImage: rockImage,
		rockPos:   image.Point{X: rockX, Y: rockY},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Clicker2")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
