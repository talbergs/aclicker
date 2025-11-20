package game_test

import (
	"os"
	"testing"

	"clicker2/game"
	"clicker2/game/events"
	"clicker2/game/eventstore"
)

func TestGameReplay(t *testing.T) {
	// Use a temporary file for the event log
	tempEventLog := "test_events.log"
	defer os.Remove(tempEventLog) // Clean up after the test

	// Initialize an event store for the original game
	originalEventStore := eventstore.NewFileEventStore(tempEventLog)

	// Create an original game instance
	originalGame := &game.Game{
		TheRock: &game.Rock{
			Health: game.InitialRockHealth,
		},
		ThePlayer: &game.Player{
			Dust:   0,
			Damage: 1,
		},
		Dispatcher: events.NewEventDispatcher(originalEventStore),
	}
	originalGame.Dispatcher.Register("DamageUpgraded", originalGame.ApplyDamageUpgradedEvent)
	originalGame.Dispatcher.Register("Click", originalGame.ApplyClickEvent)

	// Perform some actions on the original game
	originalGame.Click()
	originalGame.Click()
	originalGame.Click()
	originalGame.ThePlayer.Dust = game.UpgradeCost // Manually set dust for upgrade for testing
	originalGame.UpgradeDamage()
	originalGame.Click()
	originalGame.Click()

	// Load a new game instance from the events logged by the original game
	replayedGame, err := game.LoadGameFromEvents(originalEventStore)
	if err != nil {
		t.Fatalf("Failed to load game from events: %v", err)
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
