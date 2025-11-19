package dto

import (
	"server/common/models"
	"time"
)

type IncomingMessage struct {
	TeamID      *uint  `json:"team_id,omitempty"`
	RecipientID *uint  `json:"recipient_id,omitempty"`
	Content     string `json:"content"`
}

type MessageResponseDto struct {
	ID          uint            `json:"id"`
	SenderID    uint            `json:"sender_id"`
	Sender      UserResponseDto `json:"sender,omitempty"`
	TeamID      *uint           `json:"team_id,omitempty"`
	RecipientID *uint           `json:"recipient_id,omitempty"`
	Content     string          `json:"content"`
	CreatedAt   time.Time       `json:"created_at"`
}

func ToMessageResponseDto(msg models.Message) MessageResponseDto {
	return MessageResponseDto{
		ID:          msg.ID,
		SenderID:    msg.SenderID,
		Sender:      ToUserResponseDto(msg.Sender), // This cleans up the User object!
		TeamID:      msg.TeamID,
		RecipientID: msg.RecipientID,
		Content:     msg.Content,
		CreatedAt:   msg.CreatedAt,
	}
}
