package events

import (
	"fmt" // Added import
)

// Event is an interface for all domain events.
type Event interface {
	EventType() string
}

// DamageUpgradedEvent is dispatched when player damage is upgraded.
type DamageUpgradedEvent struct {
	PlayerID string // In a real game, this might be a player ID
	OldDamage int
	NewDamage int
	OldDust   int
	NewDust   int
}

// EventType returns the type of the DamageUpgradedEvent.
func (e *DamageUpgradedEvent) EventType() string {
	return "DamageUpgraded"
}

// UpgradePurchasedEvent is dispatched when an upgrade is purchased.
type UpgradePurchasedEvent struct {
	PlayerID string
	UpgradeID string
	NewLevel int
	OldDust int
	NewDust int
}

// EventType returns the type of the UpgradePurchasedEvent.
func (e *UpgradePurchasedEvent) EventType() string {
	return "UpgradePurchased"
}

// ClickEvent is dispatched when the rock is clicked.
type ClickEvent struct {
	PlayerID string
	DamageDealt int
	DustGained int
	RockHealthBefore int
	RockHealthAfter int
	PlayerDustBefore int
	PlayerDustAfter int
}

// EventType returns the type of the ClickEvent.
func (e *ClickEvent) EventType() string {
	return "Click"
}

// EventHandler is a function that handles a specific event.
type EventHandler func(event Event)

// EventStore defines the interface for storing events.
type EventStore interface { // Added EventStore interface
	SaveEvent(event Event) error
}

// EventDispatcher manages event handlers and dispatches events.
type EventDispatcher struct {
	handlers map[string][]EventHandler
	eventStore EventStore // Added EventStore field
}

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher(es EventStore) *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
		eventStore: es, // Can be nil for replay
	}
}

// Register registers an event handler for a specific event type.
func (ed *EventDispatcher) Register(eventType string, handler EventHandler) {
	ed.handlers[eventType] = append(ed.handlers[eventType], handler)
}

// Dispatch dispatches an event to all registered handlers for its type.
func (ed *EventDispatcher) Dispatch(event Event) {
	if ed.eventStore != nil {
		if err := ed.eventStore.SaveEvent(event); err != nil {
			// TODO: Handle error, e.g., log it
			fmt.Printf("Error saving event: %v\n", err)
		}
	}

	if handlers, ok := ed.handlers[event.EventType()]; ok {
		for _, handler := range handlers {
			handler(event)
		}
	}
}
