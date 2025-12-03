package models

import (
	"time"
)

type Message struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	SenderID uint `gorm:"not null" json:"sender_id"`
	Sender   User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`

	TeamID *uint `gorm:"index" json:"team_id,omitempty"`
	ChatID *uint `gorm:"index" json:"chat_id,omitempty"`
	Chat   *Chat `gorm:"foreignKey:ChatID" json:"chat,omitempty"`

	RecipientID *uint `gorm:"index" json:"recipient_id,omitempty"`

	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
