package hud

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// HUD represents the Heads-Up Display elements.
type HUD struct {
	// Add fields for UI elements here
	UpgradeButton image.Rectangle
}

// NewHUD creates and initializes a new HUD.
func NewHUD(screenWidth, screenHeight int) *HUD {
	buttonWidth := 150
	buttonHeight := 30
	upgradeButton := image.Rect(10, screenHeight-buttonHeight-10, 10+buttonWidth, screenHeight-10)

	return &HUD{
		UpgradeButton: upgradeButton,
	}
}

// Draw draws the HUD elements to the screen.
func (h *HUD) Draw(screen *ebiten.Image, rockHealth, playerDust, playerDamage, upgradeCost int, shadersEnabled bool) {
	// Draw the upgrade button
	ebitenutil.DrawRect(screen, float64(h.UpgradeButton.Min.X), float64(h.UpgradeButton.Min.Y), float64(h.UpgradeButton.Dx()), float64(h.UpgradeButton.Dy()), color.RGBA{R: 100, G: 100, B: 100, A: 255})

	// Draw the stats
	msg := fmt.Sprintf("Rock Health: %d\nDust: %d\nDamage: %d\nUpgrade Cost: %d\nShaders: %t (Space)", rockHealth, playerDust, playerDamage, upgradeCost, shadersEnabled)
	ebitenutil.DebugPrint(screen, msg)
}

// DrawHealthBar draws the rock's health bar.
func (h *HUD) DrawHealthBar(screen *ebiten.Image, rockPos image.Point, rockImage *ebiten.Image, currentRockHealth, initialRockHealth int) {
	barWidth := 100.0
	barHeight := 10.0
	barX := float64(rockPos.X) - (barWidth-float64(rockImage.Bounds().Dx()))/2
	barY := float64(rockPos.Y) - barHeight - 5 // 5 pixels above the rock

	healthPercentage := float64(currentRockHealth) / float64(initialRockHealth)
	if healthPercentage < 0 {
		healthPercentage = 0
	}

	// Draw health bar background
	ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Draw health bar foreground
	ebitenutil.DrawRect(screen, barX, barY, barWidth*healthPercentage, barHeight, color.RGBA{R: 0, G: 255, B: 0, A: 255})
}
