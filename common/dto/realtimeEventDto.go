package dto

import "time"

type RealtimeEventType string

const (
	RealtimeEventMessage     RealtimeEventType = "message"
	RealtimeEventTypingStart RealtimeEventType = "typing_start"
	RealtimeEventTypingStop  RealtimeEventType = "typing_stop"
)

// RealtimeEventDto is the payload sent over the WebSocket.
// It intentionally supports both message delivery and typing indicators.
type RealtimeEventDto struct {
	Type RealtimeEventType `json:"type"`

	// Routing fields (one of these should be set)
	ConversationID *uint `json:"conversation_id,omitempty"`
	TeamID         *uint `json:"team_id,omitempty"`
	RecipientID    *uint `json:"recipient_id,omitempty"`

	// Who triggered the event (sender/typer)
	UserID uint `json:"user_id"`

	// When the event happened (server time)
	Timestamp time.Time `json:"timestamp"`

	// Only set when Type == "message"
	Message *MessageResponseDto `json:"message,omitempty"`
}
