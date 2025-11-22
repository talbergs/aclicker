package game_test

import (
	"os"
	"testing"

	"clicker2/game"
	"clicker2/game/events"
	"clicker2/game/eventstore"
	"clicker2/game/errors" // Import the new errors package
)

func TestGameReplay(t *testing.T) {
	// Use a temporary file for the event log
	tempEventLog := "test_events.log"
	defer os.Remove(tempEventLog) // Clean up after the test

	// Initialize an event store for the original game
	originalEventStore := eventstore.NewFileEventStore(tempEventLog)

	// Create an original game instance with an event dispatcher that saves events
	originalGame := game.NewGame()
	originalGame.Dispatcher = events.NewEventDispatcher(originalEventStore)
	originalGame.Dispatcher.Register("Click", originalGame.ApplyClickEvent)

	// Perform some actions on the original game
	// Generate enough dust for the first stronger_pickaxe upgrade (cost 10)
	for i := 0; i < 10; i++ {
		originalGame.Click()
	}
	originalGame.PurchaseUpgrade("stronger_pickaxe")
	originalGame.Click()
	originalGame.Click()

	// Load a new game instance from the events logged by the original game
	replayedGame, err := game.LoadGameFromEvents(originalEventStore)
	if err != nil {
		t.Fatalf("Failed to load game from events: %v", err.Error())
	}

	// Assert that the state of the replayed game matches the original game
	if originalGame.TheRock.Health != replayedGame.TheRock.Health {
		t.Errorf("Rock Health mismatch: original=%d, replayed=%d", originalGame.TheRock.Health, replayedGame.TheRock.Health)
	}
	if originalGame.ThePlayer.Dust != replayedGame.ThePlayer.Dust {
		t.Errorf("Player Dust mismatch: original=%d, replayed=%d", originalGame.ThePlayer.Dust, replayedGame.ThePlayer.Dust)
	}
	if originalGame.ThePlayer.Damage != replayedGame.ThePlayer.Damage {
		t.Errorf("Player Damage mismatch: original=%d, replayed=%d", originalGame.ThePlayer.Damage, replayedGame.ThePlayer.Damage)
	}
}

func TestGameCoreMechanics(t *testing.T) {
	g := game.NewGame()

	// Test initial state
	if g.TheRock.Health != game.InitialRockHealth {
		t.Errorf("Initial rock health mismatch: got %d, want %d", g.TheRock.Health, game.InitialRockHealth)
	}
	if g.ThePlayer.Dust != 0 {
		t.Errorf("Initial player dust mismatch: got %d, want %d", g.ThePlayer.Dust, 0)
	}
	if g.ThePlayer.Damage != 1 {
		t.Errorf("Initial player damage mismatch: got %d, want %d", g.ThePlayer.Damage, 1)
	}

	// Test Click()
	initialHealth := g.TheRock.Health
	initialDust := g.ThePlayer.Dust
	g.Click()
	if g.TheRock.Health != initialHealth-g.ThePlayer.Damage {
		t.Errorf("Rock health after click mismatch: got %d, want %d", g.TheRock.Health, initialHealth-g.ThePlayer.Damage)
	}
	if g.ThePlayer.Dust != initialDust+1 {
		t.Errorf("Player dust after click mismatch: got %d, want %d", g.ThePlayer.Dust, initialDust+1)
	}

	// Test rock messages (basic check)
	if g.CurrentRockMessage == "" && len(g.RockMessages) > 0 {
		t.Errorf("Expected a rock message after click, but got empty")
	}
	if g.RockMessageTimer <= 0 {
		t.Errorf("Expected rock message timer to be set, but got %f", g.RockMessageTimer)
	}
}

func TestUpgradeSystem(t *testing.T) {
	g := game.NewGame()

	// Test "stronger_pickaxe" purchase
	g.ThePlayer.Dust = 10 // Enough for first level
	err := g.PurchaseUpgrade("stronger_pickaxe")
	if err != nil {
		t.Fatalf("Failed to purchase stronger_pickaxe: %v", err.Error())
	}
	if g.ThePlayer.Damage != 2 {
		t.Errorf("Stronger pickaxe damage mismatch: got %d, want %d", g.ThePlayer.Damage, 2)
	}
	if g.ThePlayer.Dust != 0 {
		t.Errorf("Stronger pickaxe dust mismatch: got %d, want %d", g.ThePlayer.Dust, 0)
	}

	// Test "stronger_pickaxe" max level
	g.ThePlayer.Dust = 1000 // Enough for all levels
	for i := 0; i < 4; i++ { // Purchase remaining 4 levels (total 5)
		g.PurchaseUpgrade("stronger_pickaxe")
	}
	if g.ThePlayer.Damage != 6 {
		t.Errorf("Stronger pickaxe max level damage mismatch: got %d, want %d", g.ThePlayer.Damage, 6)
	}
	err = g.PurchaseUpgrade("stronger_pickaxe")
	if err == nil || err.Code != errors.ErrUpgradeMaxLevel {
		t.Errorf("Expected max level error with code %d, got %v", errors.ErrUpgradeMaxLevel, err)
	}

	// Test "stronger_pickaxe" insufficient dust (after resetting game to ensure not max level)
	g = game.NewGame() // Reset game state
	g.ThePlayer.Dust = 0
	err = g.PurchaseUpgrade("stronger_pickaxe")
	if err == nil || err.Code != errors.ErrInsufficientDust {
		t.Errorf("Expected insufficient dust error with code %d, got %v", errors.ErrInsufficientDust, err)
	}

	// Test "auto_clicker_v0_1" purchase
	g.ThePlayer.Dust = 100 // Enough dust
	err = g.PurchaseUpgrade("auto_clicker_v0_1")
	if err != nil {
		t.Fatalf("Failed to purchase auto_clicker_v0_1: %v", err.Error())
	}
	if !g.AutoClickerActive {
		t.Errorf("Auto-clicker v0.1 not active")
	}
	if g.AutoClickerRate != 1 {
		t.Errorf("Auto-clicker v0.1 rate mismatch: got %d, want %d", g.AutoClickerRate, 1)
	}

	// Test "auto_clicker_v1_0" purchase
	g.ThePlayer.Dust = 500 // Enough dust
	err = g.PurchaseUpgrade("auto_clicker_v1_0")
	if err != nil {
		t.Fatalf("Failed to purchase auto_clicker_v1_0: %v", err.Error())
	}
	if !g.AutoClickerActive {
		t.Errorf("Auto-clicker v1.0 not active")
	}
	if g.AutoClickerRate != 5 {
		t.Errorf("Auto-clicker v1.0 rate mismatch: got %d, want %d", g.AutoClickerRate, 5)
	}
}

func TestEndings(t *testing.T) {
	g := game.NewGame()

	// Mock os.Exit to prevent test termination
	oldOsExit := game.OsExit
	game.OsExit = func(code int) {}
	defer func() { game.OsExit = oldOsExit }()

	// Purchase "heart_of_the_mountain"
	g.ThePlayer.Dust = 100000 // Enough dust
	err := g.PurchaseUpgrade("heart_of_the_mountain")
	if err != nil {
		t.Fatalf("Failed to purchase heart_of_the_mountain: %v", err.Error())
	}
	if !g.EndGameChoicePending {
		t.Errorf("EndGameChoicePending not set after purchasing heart_of_the_mountain")
	}

	// Test TakeHeart()
	g.TakeHeart()
	if g.TheRock.Health != 0 {
		t.Errorf("TakeHeart: Rock health mismatch: got %d, want %d", g.TheRock.Health, 0)
	}
	if !g.GameOver {
		t.Errorf("TakeHeart: GameOver not set")
	}

	// Reset game state for LetRest test
	g = game.NewGame()
	g.ThePlayer.Dust = 100000 // Enough dust
	g.PurchaseUpgrade("heart_of_the_mountain") // Re-purchase to set EndGameChoicePending

	// Test LetRest()
	g.LetRest()
	if !g.GameWon {
		t.Errorf("LetRest: GameWon not set")
	}
	// Further assertions for saved state would require mocking os.WriteFile
}

func TestSaveLoad(t *testing.T) {
	// Use a temporary file for saving
	tempSaveFile := "test_save.json"
	defer os.Remove(tempSaveFile)

	// Create an original game instance and modify its state
	originalGame := game.NewGame()
	originalGame.TheRock.Health = 5000000
	originalGame.ThePlayer.Dust = 12345
	originalGame.ThePlayer.Damage = 5
	err := originalGame.PurchaseUpgrade("stronger_pickaxe") // Purchase an upgrade
	if err != nil {
		t.Fatalf("Failed to purchase stronger_pickaxe in originalGame: %v", err.Error())
	}
	originalGame.AutoClickerActive = true
	originalGame.AutoClickerRate = 1
	originalGame.CurrentRockMessage = "Test Message"
	originalGame.RockMessageTimer = 1.5
	originalGame.EndGameChoicePending = true
	originalGame.GameOver = false
	originalGame.GameWon = false

	// Save the original game state
	if err := originalGame.Save(); err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	// Load a new game instance from the saved file
	loadedGame := game.NewGame() // Start with a fresh game
	if err := loadedGame.Load(); err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	// Assert that the loaded game state matches the original game state
	if originalGame.TheRock.Health != loadedGame.TheRock.Health {
		t.Errorf("Rock Health mismatch: original=%d, loaded=%d", originalGame.TheRock.Health, loadedGame.TheRock.Health)
	}
	if originalGame.ThePlayer.Dust != loadedGame.ThePlayer.Dust {
		t.Errorf("Player Dust mismatch: original=%d, loaded=%d", originalGame.ThePlayer.Dust, loadedGame.ThePlayer.Dust)
	}
	if originalGame.ThePlayer.Damage != loadedGame.ThePlayer.Damage {
		t.Errorf("Player Damage mismatch: original=%d, loaded=%d", originalGame.ThePlayer.Damage, loadedGame.ThePlayer.Damage)
	}
	if originalGame.Upgrades.PlayerUpgrades["stronger_pickaxe"] != loadedGame.Upgrades.PlayerUpgrades["stronger_pickaxe"] {
		t.Errorf("Upgrade level mismatch: original=%d, loaded=%d", originalGame.Upgrades.PlayerUpgrades["stronger_pickaxe"], loadedGame.Upgrades.PlayerUpgrades["stronger_pickaxe"])
	}
	if originalGame.AutoClickerActive != loadedGame.AutoClickerActive {
		t.Errorf("AutoClickerActive mismatch: original=%t, loaded=%t", originalGame.AutoClickerActive, loadedGame.AutoClickerActive)
	}
	if originalGame.AutoClickerRate != loadedGame.AutoClickerRate {
		t.Errorf("AutoClickerRate mismatch: original=%d, loaded=%d", originalGame.AutoClickerRate, loadedGame.AutoClickerRate)
	}
	if originalGame.CurrentRockMessage != loadedGame.CurrentRockMessage {
		t.Errorf("CurrentRockMessage mismatch: original=%s, loaded=%s", originalGame.CurrentRockMessage, loadedGame.CurrentRockMessage)
	}
	if originalGame.RockMessageTimer != loadedGame.RockMessageTimer {
		t.Errorf("RockMessageTimer mismatch: original=%f, loaded=%f", originalGame.RockMessageTimer, loadedGame.RockMessageTimer)
	}
	if originalGame.EndGameChoicePending != loadedGame.EndGameChoicePending {
		t.Errorf("EndGameChoicePending mismatch: original=%t, loaded=%t", originalGame.EndGameChoicePending, loadedGame.EndGameChoicePending)
	}
	if originalGame.GameOver != loadedGame.GameOver {
		t.Errorf("GameOver mismatch: original=%t, loaded=%t", originalGame.GameOver, loadedGame.GameOver)
	}
	if originalGame.GameWon != loadedGame.GameWon {
		t.Errorf("GameWon mismatch: original=%t, loaded=%t", originalGame.GameOver, loadedGame.GameOver)
	}
}

func TestSetStateEarlyGame(t *testing.T) {
	g := game.NewGame()
	g.SetStateEarlyGame()

	if g.TheRock.Health != game.InitialRockHealth {
		t.Errorf("EarlyGame: Rock health mismatch: got %d, want %d", g.TheRock.Health, game.InitialRockHealth)
	}
	if g.ThePlayer.Dust != 0 {
		t.Errorf("EarlyGame: Player dust mismatch: got %d, want %d", g.ThePlayer.Dust, 0)
	}
	if g.ThePlayer.Damage != 1 {
		t.Errorf("EarlyGame: Player damage mismatch: got %d, want %d", g.ThePlayer.Damage, 1)
	}
	if len(g.Upgrades.PlayerUpgrades) != 0 {
		t.Errorf("EarlyGame: Expected no upgrades, got %d", len(g.Upgrades.PlayerUpgrades))
	}
	if g.AutoClickerActive {
		t.Errorf("EarlyGame: AutoClickerActive should be false")
	}
	if g.EndGameChoicePending {
		t.Errorf("EarlyGame: EndGameChoicePending should be false")
	}
	if g.GameOver {
		t.Errorf("EarlyGame: GameOver should be false")
	}
	if g.GameWon {
		t.Errorf("EarlyGame: GameWon should be false")
	}
}

func TestSetStateMidGame(t *testing.T) {
	g := game.NewGame()
	g.SetStateMidGame()

	if g.TheRock.Health != game.InitialRockHealth/2 {
		t.Errorf("MidGame: Rock health mismatch: got %d, want %d", g.TheRock.Health, game.InitialRockHealth/2)
	}
	if g.ThePlayer.Dust != 500 {
		t.Errorf("MidGame: Player dust mismatch: got %d, want %d", g.ThePlayer.Dust, 500)
	}
	if g.ThePlayer.Damage != 5 {
		t.Errorf("MidGame: Player damage mismatch: got %d, want %d", g.ThePlayer.Damage, 5)
	}
	if g.Upgrades.PlayerUpgrades["stronger_pickaxe"] != 4 {
		t.Errorf("MidGame: Stronger pickaxe level mismatch: got %d, want %d", g.Upgrades.PlayerUpgrades["stronger_pickaxe"], 4)
	}
	if !g.AutoClickerActive {
		t.Errorf("MidGame: AutoClickerActive should be true")
	}
	if g.AutoClickerRate != 1 {
		t.Errorf("MidGame: AutoClickerRate mismatch: got %d, want %d", g.AutoClickerRate, 1)
	}
	if g.EndGameChoicePending {
		t.Errorf("MidGame: EndGameChoicePending should be false")
	}
}

func TestSetStateEndGameReady(t *testing.T) {
	g := game.NewGame()
	g.SetStateEndGameReady()

	if g.TheRock.Health != game.InitialRockHealth/10 {
		t.Errorf("EndGameReady: Rock health mismatch: got %d, want %d", g.TheRock.Health, game.InitialRockHealth/10)
	}
	if g.ThePlayer.Dust != 100000 {
		t.Errorf("EndGameReady: Player dust mismatch: got %d, want %d", g.ThePlayer.Dust, 100000)
	}
	if g.ThePlayer.Damage != 10 {
		t.Errorf("EndGameReady: Player damage mismatch: got %d, want %d", g.ThePlayer.Damage, 10)
	}
	if g.Upgrades.PlayerUpgrades["stronger_pickaxe"] != 5 {
		t.Errorf("EndGameReady: Stronger pickaxe level mismatch: got %d, want %d", g.Upgrades.PlayerUpgrades["stronger_pickaxe"], 5)
	}
	if g.Upgrades.PlayerUpgrades["auto_clicker_v1_0"] != 1 {
		t.Errorf("EndGameReady: Auto-clicker v1.0 level mismatch: got %d, want %d", g.Upgrades.PlayerUpgrades["auto_clicker_v1_0"], 1)
	}
	if !g.AutoClickerActive {
		t.Errorf("EndGameReady: AutoClickerActive should be true")
	}
	if g.AutoClickerRate != 5 {
		t.Errorf("EndGameReady: AutoClickerRate mismatch: got %d, want %d", g.AutoClickerRate, 5)
	}
	if g.EndGameChoicePending {
		t.Errorf("EndGameReady: EndGameChoicePending should be false")
	}
}
