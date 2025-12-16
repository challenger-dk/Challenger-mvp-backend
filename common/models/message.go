package models

import (
	"time"
)

type Message struct {
	ID             uint          `gorm:"primaryKey" json:"id"`
	ConversationID *uint         `gorm:"index:idx_conversation_created" json:"conversation_id,omitempty"`
	Conversation   *Conversation `gorm:"foreignKey:ConversationID" json:"conversation,omitempty"`
	SenderID       uint          `gorm:"not null" json:"sender_id"`
	Sender         User          `gorm:"foreignKey:SenderID" json:"sender,omitempty"`

	// Legacy fields for backward compatibility
	TeamID      *uint `gorm:"index" json:"team_id,omitempty"`
	RecipientID *uint `gorm:"index" json:"recipient_id,omitempty"`
	Recipient   *User `gorm:"foreignKey:RecipientID" json:"recipient,omitempty"`

	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_conversation_created" json:"created_at"`
}

