package game

import (
	"fmt"
	"clicker2/game/events"
	"clicker2/game/errors" // Import the new errors package
)

// Effect is a function that applies an upgrade to the game state.
type Effect func(g *Game)

// CostFunc is a function that calculates the cost of an upgrade, potentially based on its level.
type CostFunc func(level int) int

// Upgrade defines a single upgrade in the game.
type Upgrade struct {
	ID          string
	Name        string
	Description string
	MaxLevel    int
	Cost        CostFunc
	ApplyEffect Effect // Applied when purchased
	ReconstructEffect func(g *Game, level int) // Applied during event replay
}

// UpgradeManager manages all upgrades in the game.
type UpgradeManager struct {
	upgrades        map[string]*Upgrade
	PlayerUpgrades  map[string]int // map of upgrade ID to current level
}

// NewUpgradeManager creates a new upgrade manager and initializes the upgrades.
func NewUpgradeManager() *UpgradeManager {
	um := &UpgradeManager{
		upgrades:       make(map[string]*Upgrade),
		PlayerUpgrades: make(map[string]int),
	}
	um.Init() // Initialize upgrades when a new manager is created
	return um
}

// Init initializes or re-initializes the list of available upgrades.
func (um *UpgradeManager) Init() {
	um.upgrades = make(map[string]*Upgrade) // Clear existing upgrades before re-registering
	um.registerUpgrades()
}

func (um *UpgradeManager) registerUpgrades() {
	// Tier 1 Upgrades
	um.addUpgrade(&Upgrade{
		ID:          "stronger_pickaxe",
		Name:        "Stronger Pickaxe",
		Description: "Increases click damage by 1.",
		MaxLevel:    5,
		Cost:        func(level int) int { return 10 * (level + 1) },
		ApplyEffect: func(g *Game) {
			g.ThePlayer.Damage++
		},
		ReconstructEffect: func(g *Game, level int) {
			g.ThePlayer.Damage = 1 + level // Base damage + level
		},
	})

	um.addUpgrade(&Upgrade{
		ID:          "auto_clicker_v0_1",
		Name:        "Auto-Clicker v0.1",
		Description: "Enables a basic auto-clicker. Can be toggled.",
		MaxLevel:    1,
		Cost:        func(level int) int { return 100 },
		ApplyEffect: func(g *Game) {
			g.AutoClickerActive = true
			g.AutoClickerRate = 1
		},
		ReconstructEffect: func(g *Game, level int) {
			if level > 0 {
				g.AutoClickerActive = true
				g.AutoClickerRate = 1
			}
		},
	})

	um.addUpgrade(&Upgrade{
		ID:          "auto_clicker_v1_0",
		Name:        "Auto-Clicker v1.0",
		Description: "Upgrades the auto-clicker to be permanent and faster.",
		MaxLevel:    1,
		Cost:        func(level int) int { return 500 },
		ApplyEffect: func(g *Game) {
			g.AutoClickerActive = true
			g.AutoClickerRate = 5 // Increase rate
			// In a real game, this would also disable the toggle UI
		},
		ReconstructEffect: func(g *Game, level int) {
			if level > 0 {
				g.AutoClickerActive = true
				g.AutoClickerRate = 5
			}
		},
	})

	um.addUpgrade(&Upgrade{
		ID:          "heart_of_the_mountain",
		Name:        "The Heart of the Mountain",
		Description: "The ultimate choice. Purchase to decide the rock's fate.",
		MaxLevel:    1,
		Cost:        func(level int) int { return 100000 }, // Very high cost
		ApplyEffect: func(g *Game) {
			g.EndGameChoicePending = true
			g.CurrentRockMessage = "You have reached the Heart of the Mountain. The rock is now still. It has given all it can. You have gathered enough. Will you take the final piece, or will you let it rest?"
			g.RockMessageTimer = -1.0 // Display indefinitely until choice is made
		},
		ReconstructEffect: func(g *Game, level int) {
			if level > 0 {
				g.EndGameChoicePending = true
				g.CurrentRockMessage = "You have reached the Heart of the Mountain. The rock is now still. It has given all it can. You have gathered enough. Will you take the final piece, or will you let it rest?"
				g.RockMessageTimer = -1.0
			}
		},
	})
}

func (um *UpgradeManager) addUpgrade(u *Upgrade) {
	um.upgrades[u.ID] = u
}

// GetUpgrade returns an upgrade by its ID.
func (um *UpgradeManager) GetUpgrade(id string) (*Upgrade, *errors.GameError) {
	u, ok := um.upgrades[id]
	if !ok {
		return nil, errors.NewGameError(errors.ErrUpgradeNotFound)
	}
	return u, nil
}

// GetPlayerUpgradeLevel returns the current level of a purchased upgrade for the player.
func (um *UpgradeManager) GetPlayerUpgradeLevel(id string) int {
	level, ok := um.PlayerUpgrades[id]
	if !ok {
		return 0
	}
	return level
}

// GetAllUpgrades returns a slice of all registered upgrades.
func (um *UpgradeManager) GetAllUpgrades() []*Upgrade {
	upgrades := make([]*Upgrade, 0, len(um.upgrades))
	for _, u := range um.upgrades {
		upgrades = append(upgrades, u)
	}
	return upgrades
}

// PurchaseUpgrade attempts to purchase an upgrade for the player.
func (g *Game) PurchaseUpgrade(id string) *errors.GameError {
	u, err := g.Upgrades.GetUpgrade(id)
	if err != nil {
		return err // GetUpgrade already returns *errors.GameError
	}

	currentLevel := g.Upgrades.GetPlayerUpgradeLevel(id)
	if currentLevel >= u.MaxLevel {
		return errors.NewGameError(errors.ErrUpgradeMaxLevel)
	}

	cost := u.Cost(currentLevel)
	if g.ThePlayer.Dust < cost {
		return errors.NewGameError(errors.ErrInsufficientDust)
	}

	oldDust := g.ThePlayer.Dust // Capture old dust before deduction
	// Deduct cost and apply effect
	g.ThePlayer.Dust -= cost
	u.ApplyEffect(g)

	// Increment level
	g.Upgrades.PlayerUpgrades[id]++

	fmt.Printf("Player purchased upgrade: %s, New Level: %d\n", id, g.Upgrades.PlayerUpgrades[id])

	// Dispatch event
	g.Dispatcher.Dispatch(&events.UpgradePurchasedEvent{
		PlayerID: "player1", // Placeholder
		UpgradeID: id,
		NewLevel: g.Upgrades.PlayerUpgrades[id],
		OldDust: oldDust, // Need to capture old dust before deduction
		NewDust: g.ThePlayer.Dust,
	})
	return nil
}

