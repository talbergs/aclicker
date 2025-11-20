package game

import (
	"encoding/json"
	"os"
	"fmt"

	"clicker2/game/events"
	"clicker2/game/eventstore"
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
	es := eventstore.NewFileEventStore("events.log") // Initialize FileEventStore
	g := &Game{
		TheRock: &Rock{
			Health: InitialRockHealth,
		},
		ThePlayer: &Player{
			Dust:   0,
			Damage: 1,
		},
		Dispatcher: events.NewEventDispatcher(es), // Pass EventStore to dispatcher
	}
	g.Dispatcher.Register("DamageUpgraded", g.ApplyDamageUpgradedEvent)
	g.Dispatcher.Register("Click", g.ApplyClickEvent)
	return g
}

// Click handles the logic for a single click on the rock.
func (g *Game) Click() {
	rockHealthBefore := g.TheRock.Health
	playerDustBefore := g.ThePlayer.Dust

	damageDealt := g.ThePlayer.Damage
	dustGained := 1

	// Dispatch event
	g.Dispatcher.Dispatch(&events.ClickEvent{
		PlayerID: "player1", // Placeholder
		DamageDealt: damageDealt,
		DustGained: dustGained,
		RockHealthBefore: rockHealthBefore,
		RockHealthAfter: rockHealthBefore - damageDealt,
		PlayerDustBefore: playerDustBefore,
		PlayerDustAfter: playerDustBefore + dustGained,
	})
}

// ApplyClickEvent applies the state changes from a ClickEvent.
func (g *Game) ApplyClickEvent(event events.Event) {
	if e, ok := event.(*events.ClickEvent); ok {
		g.TheRock.Health = e.RockHealthAfter
		g.ThePlayer.Dust = e.PlayerDustAfter
	}
}

// ReplayEvents takes a slice of events and dispatches them to reconstruct the game state.
func (g *Game) ReplayEvents(evs []events.Event) {
	for _, event := range evs {
		g.Dispatcher.Dispatch(event)
	}
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

// LoadGameFromEvents creates a new game instance and replays events from the provided EventStore.
func LoadGameFromEvents(es eventstore.EventStore) (*Game, error) {
	// Create a new game instance with an event dispatcher that does NOT save events during replay
	g := &Game{
		TheRock: &Rock{
			Health: InitialRockHealth,
		},
		ThePlayer: &Player{
			Dust:   0,
			Damage: 1,
		},
		Dispatcher: events.NewEventDispatcher(nil), // Pass nil for EventStore during replay
	}
	g.Dispatcher.Register("DamageUpgraded", g.ApplyDamageUpgradedEvent)
	g.Dispatcher.Register("Click", g.ApplyClickEvent)

	// Load all events from the event store
	loadedEvents, err := es.LoadEvents()
	if err != nil {
		return nil, fmt.Errorf("failed to load events from event store: %w", err)
	}

	// Replay the events to reconstruct the game state
	g.ReplayEvents(loadedEvents)

	// After replay, set up the dispatcher to save new events
	g.Dispatcher = events.NewEventDispatcher(es) // Re-initialize with the actual EventStore
	g.Dispatcher.Register("DamageUpgraded", g.ApplyDamageUpgradedEvent)
	g.Dispatcher.Register("Click", g.ApplyClickEvent)


	return g, nil
}

