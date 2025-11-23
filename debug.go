package main

import (
	"clicker2/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Debug holds debug-related state and functions.
type Debug struct {
	AutoClickerEnabled bool
	AutoClickerSpeed   int
	gameState          *game.Game // Reference to the main game state
}

// NewDebug creates a new Debug instance.
func NewDebug(gs *game.Game) *Debug {
	return &Debug{
		AutoClickerEnabled: false,
		AutoClickerSpeed:   0,
		gameState:          gs,
	}
}

// HandleDeveloperKeybinds processes developer-specific key inputs.
func (d *Debug) HandleDeveloperKeybinds() {
	// Developer: Toggle debug auto-clicker
	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		d.AutoClickerEnabled = !d.AutoClickerEnabled
		if d.AutoClickerEnabled {
			log.Println("Debug Auto-Clicker ENABLED")
		} else {
			log.Println("Debug Auto-Clicker DISABLED")
		}
	}

	// Developer: Cycle debug auto-clicker speed
	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		if d.AutoClickerEnabled {
			switch d.AutoClickerSpeed {
			case 0:
				d.AutoClickerSpeed = 1
			case 1:
				d.AutoClickerSpeed = 5
			case 5:
				d.AutoClickerSpeed = 10
			case 10:
				d.AutoClickerSpeed = 0 // Turn off debug auto-clicker
				d.AutoClickerEnabled = false
			}
			log.Printf("Debug Auto-Clicker Speed: %d clicks/second", d.AutoClickerSpeed)
		} else {
			log.Println("Debug Auto-Clicker is disabled. Press F4 to enable.")
		}
	}

	// Developer: Load specific game states
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		d.gameState.SetStateEarlyGame()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		d.gameState.SetStateMidGame()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		d.gameState.SetStateEndGameReady()
	}
}

// UpdateDebugAutoClicker handles the logic for the debug auto-clicker.
func (d *Debug) UpdateDebugAutoClicker() {
	if d.AutoClickerEnabled && d.AutoClickerSpeed > 0 {
		// Calculate frames per click
		framesPerClick := ebiten.TPS() / d.AutoClickerSpeed
		if framesPerClick <= 0 { // Ensure at least one click per frame if speed is very high
			framesPerClick = 1
		}
		if int(ebiten.ActualTPS())%framesPerClick == 0 {
			d.gameState.Click()
		}
	}
}
