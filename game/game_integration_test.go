package game_test

import (
	"clicker2/game"
	"os"
	"testing"
)

// TestGameLifecycleAndExit tests the basic game lifecycle, including clicks,
// state changes, and the new graceful exit mechanism.
func TestGameLifecycleAndExit(t *testing.T) {
	// Mock os.Exit to prevent actual program termination during tests
	originalOsExit := game.OsExit
	game.OsExit = func(code int) {
		// Do nothing, or log for debugging if needed
		t.Logf("os.Exit called with code %d (mocked)", code)
	}
	defer func() {
		game.OsExit = originalOsExit // Restore original os.Exit after test
	}()

	// 1. Initialize a new game
	g := game.NewGame()

	// Assert initial state
	if g.TheRock.Health != game.InitialRockHealth {
		t.Errorf("Expected initial rock health %d, got %d", game.InitialRockHealth, g.TheRock.Health)
	}
	if g.ThePlayer.Dust != 0 {
		t.Errorf("Expected initial player dust 0, got %d", g.ThePlayer.Dust)
	}
	if g.ShouldExit {
		t.Error("Expected ShouldExit to be false initially")
	}

	// 2. Simulate clicks and observe state changes
	initialHealth := g.TheRock.Health
	initialDust := g.ThePlayer.Dust

	clicks := 10
	for i := 0; i < clicks; i++ {
		g.Click()
	}

	expectedHealth := initialHealth - (clicks * g.ThePlayer.Damage)
	expectedDust := initialDust + clicks

	if g.TheRock.Health != expectedHealth {
		t.Errorf("After %d clicks, expected rock health %d, got %d", clicks, expectedHealth, g.TheRock.Health)
	}
	if g.ThePlayer.Dust != expectedDust {
		t.Errorf("After %d clicks, expected player dust %d, got %d", clicks, expectedDust, g.ThePlayer.Dust)
	}

	// 3. Simulate game ending (e.g., by taking the heart)
	g.EndGameChoicePending = true // Simulate the condition for end-game choice
	g.TakeHeart()

	// Assert that ShouldExit is true
	if !g.ShouldExit {
		t.Error("Expected ShouldExit to be true after TakeHeart()")
	}
	if !g.GameOver {
		t.Error("Expected GameOver to be true after TakeHeart()")
	}

	// 4. Test Save/Load functionality
	// Reset ShouldExit for saving purposes, as a real game might save before exiting
	g.ShouldExit = false

	// Create a temporary save file
	tempSaveFile := "test_save.json"
	defer os.Remove(tempSaveFile) // Clean up temp file

	// Save the current game state
	if err := g.SaveToFile(tempSaveFile); err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	// Create a new game instance and load the saved state
	loadedGame := game.NewGame()
	if err := loadedGame.LoadFromFile(tempSaveFile); err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	// Assert that the loaded game state matches the saved state
	if loadedGame.TheRock.Health != g.TheRock.Health {
		t.Errorf("Loaded rock health %d does not match saved %d", loadedGame.TheRock.Health, g.TheRock.Health)
	}
	if loadedGame.ThePlayer.Dust != g.ThePlayer.Dust {
		t.Errorf("Loaded player dust %d does not match saved %d", loadedGame.ThePlayer.Dust, g.ThePlayer.Dust)
	}
	if loadedGame.GameOver != g.GameOver {
		t.Errorf("Loaded GameOver %t does not match saved %t", loadedGame.GameOver, g.GameOver)
	}
	if loadedGame.GameWon != g.GameWon {
		t.Errorf("Loaded GameWon %t does not match saved %t", loadedGame.GameWon, g.GameWon)
	}
	// Add more assertions for other relevant fields if necessary
}
