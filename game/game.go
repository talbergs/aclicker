package game

import (
	"encoding/json"
	"os"
)

const (
	InitialRockHealth = 10000000
	UpgradeCost       = 10
	saveFile          = "save.json"
)

// Rock represents the entity that is clicked.
type Rock struct {
	Health int
}

// Player represents the user's state.
type Player struct {
	Dust   int
	Damage int
}

// Game holds the overall game state.
type Game struct {
	TheRock   *Rock
	ThePlayer *Player
}

// NewGame creates a new game state with initial values.
func NewGame() *Game {
	return &Game{
		TheRock: &Rock{
			Health: InitialRockHealth,
		},
		ThePlayer: &Player{
			Dust:   0,
			Damage: 1,
		},
	}
}

// Click handles the logic for a single click on the rock.
func (g *Game) Click() {
	g.TheRock.Health -= g.ThePlayer.Damage
	g.ThePlayer.Dust++
}

// UpgradeDamage increases the player's damage if they have enough dust.
func (g *Game) UpgradeDamage() {
	if g.ThePlayer.Dust >= UpgradeCost {
		g.ThePlayer.Dust -= UpgradeCost
		g.ThePlayer.Damage++
	}
}

// Save serializes the game state to a file.
func (g *Game) Save() error {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(saveFile, data, 0644)
}

// Load deserializes the game state from a file.
func (g *Game) Load() error {
	data, err := os.ReadFile(saveFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, g)
}

