package game

import (
	"encoding/json"
	"os"
	"fmt"
	"math/rand" // Added for random message selection
	"log" // Added for logging game endings

	"clicker2/game/events"
	"clicker2/game/eventstore"
	"clicker2/game/errors" // Import the new errors package
)

// osExit is a package-level variable that can be mocked for testing os.Exit
var OsExit = os.Exit

const InitialRockHealth = 10000000
var SaveFile = "save.json" // Exported for testing

// Save serializes the game state to a file.
func (g *Game) Save() error {
	return g.SaveToFile(SaveFile)
}

// SaveToFile serializes the game state to the specified file path.
func (g *Game) SaveToFile(path string) error {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

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
	Upgrades  *UpgradeManager
	Dispatcher *events.EventDispatcher
	AutoClickerActive bool
	AutoClickerRate   int
	CurrentRockMessage string
	RockMessageTimer   float64 // Duration for which the message is displayed
	RockMessages       []string
	EndGameChoicePending bool
	GameOver             bool
	GameWon              bool
	ShouldExit           bool // New field to signal game termination
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
		Upgrades:   NewUpgradeManager(),
		Dispatcher: events.NewEventDispatcher(es), // Pass EventStore to dispatcher
		AutoClickerActive: false,
		AutoClickerRate:   0,
		CurrentRockMessage: "",
		RockMessageTimer:   0.0,
		RockMessages: []string{
			"A gentle hum emanates from within...",
			"You feel a faint tremor.",
			"The rock seems... content.",
			"A tiny shard breaks off, almost imperceptibly.",
			"You hear a soft, distant sigh.",
		},
		EndGameChoicePending: false,
		GameOver:             false,
		GameWon:              false,
		ShouldExit:           false, // Initialize ShouldExit to false
	}
	g.Dispatcher.Register("Click", g.ApplyClickEvent)
	g.Dispatcher.Register("UpgradePurchased", g.ApplyUpgradePurchasedEvent)
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

	// Trigger a rock message
	if len(g.RockMessages) > 0 {
		randomIndex := rand.Intn(len(g.RockMessages))
		g.CurrentRockMessage = g.RockMessages[randomIndex]
		g.RockMessageTimer = 3.0 // Display message for 3 seconds
	}
}

// ApplyClickEvent applies the state changes from a ClickEvent.
func (g *Game) ApplyClickEvent(event events.Event) {
	if e, ok := event.(*events.ClickEvent); ok {
		g.TheRock.Health = e.RockHealthAfter
		g.ThePlayer.Dust = e.PlayerDustAfter
	}
}

// ApplyUpgradePurchasedEvent applies the state changes from an UpgradePurchasedEvent.
func (g *Game) ApplyUpgradePurchasedEvent(event events.Event) {
	if e, ok := event.(*events.UpgradePurchasedEvent); ok {
		g.Upgrades.PlayerUpgrades[e.UpgradeID] = e.NewLevel
		g.ThePlayer.Dust = e.NewDust
		// Apply the effect of the upgrade again to ensure state consistency
		// This is important for replay, as the effect might modify other game state
		// that isn't directly part of the event (e.g., g.ThePlayer.Damage)
		upgrade, err := g.Upgrades.GetUpgrade(e.UpgradeID)
		if err != nil { // err is now *errors.GameError
			log.Printf("Error getting upgrade %s during replay: %v", e.UpgradeID, err.Error())
			return
		}
		upgrade.ReconstructEffect(g, e.NewLevel)
	}
}

// PurchaseUpgrade handles the logic for purchasing an upgrade.
func (g *Game) PurchaseUpgrade(upgradeID string) *errors.GameError {
	upgrade, err := g.Upgrades.GetUpgrade(upgradeID)
	if err != nil {
		return err
	}

	if g.Upgrades.PlayerUpgrades[upgradeID] >= upgrade.MaxLevel {
		return errors.NewGameError(errors.ErrUpgradeMaxLevel, "upgrade at max level")
	}

	cost := upgrade.Cost(g.Upgrades.PlayerUpgrades[upgradeID])
	if g.ThePlayer.Dust < cost {
		return errors.NewGameError(errors.ErrInsufficientDust, "not enough dust to purchase upgrade")
	}

	g.ThePlayer.Dust -= cost
	g.Upgrades.PlayerUpgrades[upgradeID]++
	upgrade.ApplyEffect(g)

	g.Dispatcher.Dispatch(&events.UpgradePurchasedEvent{
		UpgradeID: upgradeID,
		NewLevel:  g.Upgrades.PlayerUpgrades[upgradeID],
		NewDust:   g.ThePlayer.Dust,
	})

	return nil
}

// ReplayEvents takes a slice of events and dispatches them to reconstruct the game state.
func (g *Game) ReplayEvents(evs []events.Event) {
	for _, event := range evs {
		g.Dispatcher.Dispatch(event)
	}
}





// Load deserializes the game state from a file.
func (g *Game) Load() error {
	return g.LoadFromFile(SaveFile)
}

// LoadFromFile deserializes the game state from the specified file path.
func (g *Game) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, g); err != nil {
		return err
	}
	g.Upgrades.Init() // Re-initialize the upgrades map after loading
	return nil
}

// LoadGameFromEvents creates a new game instance and replays events from the provided EventStore.
func LoadGameFromEvents(es eventstore.EventStore) (*Game, *errors.GameError) {
	// Create a new game instance with an event dispatcher that does NOT save events during replay
	g := &Game{
		TheRock: &Rock{
			Health: InitialRockHealth,
		},
		ThePlayer: &Player{
			Dust:   0,
			Damage: 1,
		},
		Upgrades:   NewUpgradeManager(), // Initialize UpgradeManager
		Dispatcher: events.NewEventDispatcher(nil), // Pass nil for EventStore during replay
		AutoClickerActive: false,
		AutoClickerRate:   0,
		CurrentRockMessage: "",
		RockMessageTimer:   0.0,
		RockMessages: []string{
			"A gentle hum emanates from within...",
			"You feel a faint tremor.",
			"The rock seems... content.",
			"A tiny shard breaks off, almost imperceptibly.",
			"You hear a soft, distant sigh.",
		},
		EndGameChoicePending: false,
		GameOver:             false,
		GameWon:              false,
		ShouldExit:           false, // Initialize ShouldExit to false
	}
	g.Upgrades.Init() // Initialize upgrades
	g.Dispatcher.Register("Click", g.ApplyClickEvent)
	g.Dispatcher.Register("UpgradePurchased", g.ApplyUpgradePurchasedEvent)

	// Load all events from the event store
	loadedEvents, err := es.LoadEvents()
	if err != nil {
		return nil, errors.NewGameError(errors.ErrUnknown, fmt.Sprintf("failed to load events from event store: %v", err))
	}

	// Replay the events to reconstruct the game state
	g.ReplayEvents(loadedEvents)

	// After replay, set up the dispatcher to save new events
	g.Dispatcher = events.NewEventDispatcher(es) // Re-initialize with the actual EventStore
	g.Dispatcher.Register("Click", g.ApplyClickEvent)
	g.Dispatcher.Register("UpgradePurchased", g.ApplyUpgradePurchasedEvent)


	return g, nil
}

// TakeHeart implements the "Bad Ending" logic.
func (g *Game) TakeHeart() {
	log.Println("Bad Ending: You took the Heart of the Mountain.")
	g.TheRock.Health = 0 // Shatter the rock
	g.CurrentRockMessage = "The mountain is no more. You are alone with your dust."
	g.RockMessageTimer = -1.0 // Display indefinitely
	g.GameOver = true
	// In a real game, you might show a final screen before exiting.
	g.ShouldExit = true // Signal main loop to terminate
}

// LetRest implements the "Good Ending" logic.
func (g *Game) LetRest() {
	log.Println("Good Ending: You let the Heart of the Mountain rest.")
	g.CurrentRockMessage = "The rock is at peace. You have won."
	g.RockMessageTimer = -1.0 // Display indefinitely
	g.GameWon = true
	// Save the game in its "won" state
	if err := g.Save(); err != nil {
		log.Printf("Error saving game after winning: %v", err)
	}
	// In a real game, you might show a final screen before exiting.
	g.ShouldExit = true // Signal main loop to terminate
}

// SetStateEarlyGame sets the game state to an early game scenario.
func (g *Game) SetStateEarlyGame() {
	g.TheRock.Health = InitialRockHealth
	g.ThePlayer.Dust = 0
	g.ThePlayer.Damage = 1
	g.Upgrades.PlayerUpgrades = make(map[string]int) // Clear upgrades
	g.Upgrades.Init() // Re-initialize upgrade definitions
	g.AutoClickerActive = false
	g.AutoClickerRate = 0
	g.CurrentRockMessage = ""
	g.RockMessageTimer = 0.0
	g.EndGameChoicePending = false
	g.GameOver = false
	g.GameWon = false
	log.Println("Game state set to Early Game.")
}

// SetStateMidGame sets the game state to a mid-game scenario.
func (g *Game) SetStateMidGame() {
	g.SetStateEarlyGame() // Start from early game state
	g.TheRock.Health = InitialRockHealth / 2
	g.ThePlayer.Dust = 500
	g.ThePlayer.Damage = 5
	g.Upgrades.PlayerUpgrades["stronger_pickaxe"] = 4 // Some upgrades
	g.Upgrades.PlayerUpgrades["auto_clicker_v0_1"] = 1
	g.AutoClickerActive = true
	g.AutoClickerRate = 1
	log.Println("Game state set to Mid Game.")
}

// SetStateEndGameReady sets the game state to be ready for the end-game choice.
func (g *Game) SetStateEndGameReady() {
	g.SetStateMidGame() // Start from mid game state
	g.TheRock.Health = InitialRockHealth / 10
	g.ThePlayer.Dust = 100000 // Enough to buy Heart of the Mountain
	g.ThePlayer.Damage = 10
	g.Upgrades.PlayerUpgrades["stronger_pickaxe"] = 5 // Max stronger pickaxe
	g.Upgrades.PlayerUpgrades["auto_clicker_v1_0"] = 1 // Permanent auto-clicker
	g.AutoClickerActive = true
	g.AutoClickerRate = 5
	log.Println("Game state set to End Game Ready.")
}

