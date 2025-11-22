package errors

// ErrorCode represents a specific error condition in the game.
type ErrorCode int

const (
	// General errors
	ErrUnknown ErrorCode = iota + 1000 // Start from a base to avoid conflicts

	// Upgrade-related errors
	ErrInsufficientDust
	ErrUpgradeMaxLevel
	ErrUpgradeNotFound

	// Event-related errors
	ErrUnknownEventType
)

// errorMessages maps ErrorCode to a default English message.
// In a full i18n system, this would be loaded from locale files.
var errorMessages = map[ErrorCode]string{
	ErrUnknown:          "An unknown error occurred.",
	ErrInsufficientDust: "Not enough dust to purchase upgrade.",
	ErrUpgradeMaxLevel:  "Upgrade already at max level.",
	ErrUpgradeNotFound:  "Upgrade not found.",
	ErrUnknownEventType: "Unknown event type encountered.",
}

// GetErrorMessage returns the human-readable message for a given ErrorCode.
// If the code is unknown, it returns a generic unknown error message.
func GetErrorMessage(code ErrorCode) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return errorMessages[ErrUnknown]
}

// GameError is a custom error type that includes an ErrorCode.
type GameError struct {
	Code    ErrorCode
	Message string
}

func (e *GameError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return GetErrorMessage(e.Code)
}

// NewGameError creates a new GameError with a given code and an optional custom message.
func NewGameError(code ErrorCode, msg ...string) *GameError {
	ge := &GameError{Code: code}
	if len(msg) > 0 {
		ge.Message = msg[0]
	}
	return ge
}