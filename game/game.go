package game

import (
	"encoding/json"
	"os"

	"clicker2/game/events"
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
	Dispatcher *events.EventDispatcher
}

// NewGame creates a new game state with initial values.
func NewGame() *Game {
	g := &Game{
		TheRock: &Rock{
			Health: InitialRockHealth,
		},
		ThePlayer: &Player{
			Dust:   0,
			Damage: 1,
		},
		Dispatcher: events.NewEventDispatcher(),
	}
	g.Dispatcher.Register("DamageUpgraded", g.ApplyDamageUpgradedEvent)
	return g
}

// Click handles the logic for a single click on the rock.
func (g *Game) Click() {
	g.TheRock.Health -= g.ThePlayer.Damage
	g.ThePlayer.Dust++
}

// ApplyDamageUpgradedEvent applies the state changes from a DamageUpgradedEvent.
func (g *Game) ApplyDamageUpgradedEvent(event events.Event) {
	if e, ok := event.(*events.DamageUpgradedEvent); ok {
		g.ThePlayer.Dust = e.NewDust
		g.ThePlayer.Damage = e.NewDamage
	}
}

// UpgradeDamage dispatches an event to upgrade the player's damage if they have enough dust.
func (g *Game) UpgradeDamage() {
	if g.ThePlayer.Dust >= UpgradeCost {
		oldDust := g.ThePlayer.Dust
		oldDamage := g.ThePlayer.Damage

		// Calculate new state
		newDust := oldDust - UpgradeCost
		newDamage := oldDamage + 1

		// Dispatch event
		g.Dispatcher.Dispatch(&events.DamageUpgradedEvent{
			PlayerID: "player1", // Placeholder, in a real game this would be dynamic
			OldDamage: oldDamage,
			NewDamage: newDamage,
			OldDust: oldDust,
			NewDust: newDust,
		})
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

