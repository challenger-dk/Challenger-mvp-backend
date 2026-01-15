package dto

import (
	"server/common/models"
	"time"
)

type IncomingMessage struct {
	Type string `json:"type,omitempty"`

	ConversationID *uint `json:"conversation_id,omitempty"`
	TeamID         *uint `json:"team_id,omitempty"`      // Legacy team messaging
	RecipientID    *uint `json:"recipient_id,omitempty"` // Legacy direct messaging

	// Only used for Type == "message"
	Content string `json:"content" validate:"sanitize"`
}

type MessageResponseDto struct {
	ID             uint            `json:"id"`
	ConversationID *uint           `json:"conversation_id,omitempty"`
	SenderID       uint            `json:"sender_id"`
	Sender         UserResponseDto `json:"sender,omitempty"`
	TeamID         *uint           `json:"team_id,omitempty"`
	RecipientID    *uint           `json:"recipient_id,omitempty"`
	Content        string          `json:"content"`
	CreatedAt      time.Time       `json:"created_at"`
}

func ToMessageResponseDto(msg models.Message) MessageResponseDto {
	return MessageResponseDto{
		ID:             msg.ID,
		ConversationID: msg.ConversationID,
		SenderID:       msg.SenderID,
		Sender:         ToUserResponseDto(msg.Sender),
		TeamID:         msg.TeamID,
		RecipientID:    msg.RecipientID,
		Content:        msg.Content,
		CreatedAt:      msg.CreatedAt,
	}
}
