package eventstore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"clicker2/game/events"
)

// EventStore defines the interface for storing and loading events.
type EventStore interface {
	SaveEvent(event events.Event) error
	LoadEvents() ([]events.Event, error)
}

// FileEventStore implements EventStore for file-based persistence.
type FileEventStore struct {
	filePath string
	mu       sync.Mutex // Protects file writes
}

// NewFileEventStore creates a new FileEventStore.
func NewFileEventStore(filePath string) *FileEventStore {
	return &FileEventStore{
		filePath: filePath,
	}
}

// SaveEvent serializes an event and appends it to the file.
func (fs *FileEventStore) SaveEvent(event events.Event) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := os.OpenFile(fs.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open event store file: %w", err)
	}
	defer file.Close()

	// Wrap the event with its type for deserialization
	eventWrapper := struct {
		Type string        `json:"type"`
		Data json.RawMessage `json:"data"`
	}{
		Type: event.EventType(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	eventWrapper.Data = data

	wrappedData, err := json.Marshal(eventWrapper)
	if err != nil {
		return fmt.Errorf("failed to marshal event wrapper: %w", err)
	}

	_, err = file.Write(append(wrappedData, '\n'))
	if err != nil {
		return fmt.Errorf("failed to write event to file: %w", err)
	}
	return nil
}

// LoadEvents reads all events from the file and deserializes them.
func (fs *FileEventStore) LoadEvents() ([]events.Event, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, err := os.OpenFile(fs.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open event store file: %w", err)
	}
	defer file.Close()

	var loadedEvents []events.Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var eventWrapper struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(line, &eventWrapper); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event wrapper: %w", err)
		}

		var event events.Event
		switch eventWrapper.Type {
		case "DamageUpgraded":
			var e events.DamageUpgradedEvent
			if err := json.Unmarshal(eventWrapper.Data, &e); err != nil {
				return nil, fmt.Errorf("failed to unmarshal DamageUpgradedEvent: %w", err)
			}
			event = &e
		case "Click": // Added case for ClickEvent
			var e events.ClickEvent
			if err := json.Unmarshal(eventWrapper.Data, &e); err != nil {
				return nil, fmt.Errorf("failed to unmarshal ClickEvent: %w", err)
			}
			event = &e
		case "UpgradePurchased": // Added case for UpgradePurchasedEvent
			var e events.UpgradePurchasedEvent
			if err := json.Unmarshal(eventWrapper.Data, &e); err != nil {
				return nil, fmt.Errorf("failed to unmarshal UpgradePurchasedEvent: %w", err)
			}
			event = &e
		// Add other event types here as they are defined
		default:
			return nil, fmt.Errorf("unknown event type: %s", eventWrapper.Type)
		}
		loadedEvents = append(loadedEvents, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading event store file: %w", err)
	}

	return loadedEvents, nil
}
