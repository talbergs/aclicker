package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Game implements ebiten.Game interface.
type Game struct{}

// Update proceeds the game state.
// Update is called every tick (1/60 second).
func (g *Game) Update() error {
	// Write your game's logical update here.
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60 second).
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering commands here.
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return the screen size.
// For more detailed explanations, see https://github.com/hajimehoshi/ebiten/v2/wiki/Ebiten's-viewports.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Clicker2")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
