package hud

import (
	"fmt"
	"image"
	"image/color"

	"clicker2/game" // Import the game package
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// UpgradeButton represents a clickable UI element for an upgrade.
type UpgradeButton struct {
	Bounds    image.Rectangle
	UpgradeID string
}

// HUD represents the Heads-Up Display elements.
type HUD struct {
	UpgradeButtons []UpgradeButton
	TakeHeartButton image.Rectangle
	LetRestButton   image.Rectangle
}

// NewHUD creates and initializes a new HUD.
func NewHUD(screenWidth, screenHeight int, um *game.UpgradeManager) *HUD {
	hud := &HUD{
		UpgradeButtons: []UpgradeButton{},
	}

	buttonWidth := 150
	buttonHeight := 30
	padding := 10
	currentY := screenHeight - padding - buttonHeight

	// Get all upgrades from the manager
	allUpgrades := um.GetAllUpgrades() // Assuming GetAllUpgrades method exists or will be added

	for _, upgrade := range allUpgrades {
		buttonRect := image.Rect(padding, currentY, padding+buttonWidth, currentY+buttonHeight)
		hud.UpgradeButtons = append(hud.UpgradeButtons, UpgradeButton{
			Bounds:    buttonRect,
			UpgradeID: upgrade.ID,
		})
		currentY -= (buttonHeight + padding) // Move up for the next button
	}

	// Initialize end-game choice buttons (initially hidden)
	choiceButtonWidth := 200
	choiceButtonHeight := 50
	choiceButtonX := screenWidth/2 - choiceButtonWidth/2
	choiceButtonY := screenHeight/2 + 50 // Below the rock message

	hud.TakeHeartButton = image.Rect(choiceButtonX, choiceButtonY, choiceButtonX+choiceButtonWidth, choiceButtonY+choiceButtonHeight)
	hud.LetRestButton = image.Rect(choiceButtonX, choiceButtonY+choiceButtonHeight+padding, choiceButtonX+choiceButtonWidth, choiceButtonY+choiceButtonHeight*2+padding)

	return hud
}

// Draw draws the HUD elements to the screen.
func (h *HUD) Draw(screen *ebiten.Image, g *game.Game, shadersEnabled bool) {
	// Draw the stats
	msg := fmt.Sprintf("Rock Health: %d\nDust: %d\nDamage: %d\nShaders: %t (Space)", g.TheRock.Health, g.ThePlayer.Dust, g.ThePlayer.Damage, shadersEnabled)
	ebitenutil.DebugPrint(screen, msg)

	// Draw rock message if active
	if g.RockMessageTimer > 0 && g.CurrentRockMessage != "" {
		messageX := screen.Bounds().Dx()/2 - 100 // Center the message
		messageY := screen.Bounds().Dy()/2 + 100 // Below the rock
		ebitenutil.DebugPrintAt(screen, g.CurrentRockMessage, messageX, messageY)
	}

	// Draw upgrade buttons
	for _, btn := range h.UpgradeButtons {
		// Draw button background
		ebitenutil.DrawRect(screen, float64(btn.Bounds.Min.X), float64(btn.Bounds.Min.Y), float64(btn.Bounds.Dx()), float64(btn.Bounds.Dy()), color.RGBA{R: 100, G: 100, B: 100, A: 255})

		// Get upgrade details
		upgrade, err := g.Upgrades.GetUpgrade(btn.UpgradeID)
		if err != nil {
			// Log error or draw placeholder
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Error: %s", err.Error()), btn.Bounds.Min.X+5, btn.Bounds.Min.Y+5)
			continue
		}

		currentLevel := g.Upgrades.GetPlayerUpgradeLevel(btn.UpgradeID)
		cost := upgrade.Cost(currentLevel)

		// Draw upgrade text
		upgradeText := fmt.Sprintf("%s\nLvl: %d Cost: %d", upgrade.Name, currentLevel, cost)
		ebitenutil.DebugPrintAt(screen, upgradeText, btn.Bounds.Min.X+5, btn.Bounds.Min.Y+5)
	}

	// Draw end-game choice buttons if pending
	if g.EndGameChoicePending {
		// Draw "Take the Heart" button
		ebitenutil.DrawRect(screen, float64(h.TakeHeartButton.Min.X), float64(h.TakeHeartButton.Min.Y), float64(h.TakeHeartButton.Dx()), float64(h.TakeHeartButton.Dy()), color.RGBA{R: 200, G: 50, B: 50, A: 255})
		ebitenutil.DebugPrintAt(screen, "Take the Heart", h.TakeHeartButton.Min.X+10, h.TakeHeartButton.Min.Y+15)

		// Draw "Let it Rest" button
		ebitenutil.DrawRect(screen, float64(h.LetRestButton.Min.X), float64(h.LetRestButton.Min.Y), float64(h.LetRestButton.Dx()), float64(h.LetRestButton.Dy()), color.RGBA{R: 50, G: 200, B: 50, A: 255})
		ebitenutil.DebugPrintAt(screen, "Let it Rest", h.LetRestButton.Min.X+10, h.LetRestButton.Min.Y+15)
	}
}

// GetClickedUpgradeID checks if a given point (e.g., mouse click) is within any upgrade button's bounds.
// Returns the UpgradeID of the clicked button, or an empty string if no button was clicked.
func (h *HUD) GetClickedUpgradeID(cursorPoint image.Point) string {
	for _, btn := range h.UpgradeButtons {
		if cursorPoint.In(btn.Bounds) {
			return btn.UpgradeID
		}
	}
	return ""
}

// GetClickedChoiceID checks if a given point (e.g., mouse click) is within any end-game choice button's bounds.
// Returns the ID of the clicked choice ("take_heart", "let_rest"), or an empty string if no choice button was clicked.
func (h *HUD) GetClickedChoiceID(cursorPoint image.Point) string {
	if cursorPoint.In(h.TakeHeartButton) {
		return "take_heart"
	}
	if cursorPoint.In(h.LetRestButton) {
		return "let_rest"
	}
	return ""
}

// DrawHealthBar draws the rock's health bar.
func (h *HUD) DrawHealthBar(screen *ebiten.Image, rockPos image.Point, rockImage *ebiten.Image, currentRockHealth, initialRockHealth int) {
	barWidth := 100.0
	barHeight := 10.0
	barX := float64(rockPos.X) / 2
	barY := float64(rockPos.Y) / 2

	healthPercentage := float64(currentRockHealth) / float64(initialRockHealth)
	if healthPercentage < 0 {
		healthPercentage = 0
	}

	// Draw health bar background
	ebitenutil.DrawRect(screen, barX, barY, barWidth, barHeight, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Draw health bar foreground
	ebitenutil.DrawRect(screen, barX, barY, barWidth*healthPercentage, barHeight, color.RGBA{R: 0, G: 255, B: 0, A: 255})
}
