package events

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

// EventHandler is a function that handles a specific event.
type EventHandler func(event Event)

// EventDispatcher manages event handlers and dispatches events.
type EventDispatcher struct {
	handlers map[string][]EventHandler
}

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

// Register registers an event handler for a specific event type.
func (ed *EventDispatcher) Register(eventType string, handler EventHandler) {
	ed.handlers[eventType] = append(ed.handlers[eventType], handler)
}

// Dispatch dispatches an event to all registered handlers for its type.
func (ed *EventDispatcher) Dispatch(event Event) {
	if handlers, ok := ed.handlers[event.EventType()]; ok {
		for _, handler := range handlers {
			handler(event)
		}
	}
}
