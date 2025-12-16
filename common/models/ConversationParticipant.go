package models

import (
	"time"
)

type ConversationParticipant struct {
	ConversationID uint         `gorm:"primaryKey;autoIncrement:false" json:"conversation_id"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"conversation,omitempty"`
	UserID         uint         `gorm:"primaryKey;autoIncrement:false" json:"user_id"`
	User           User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user,omitempty"`
	JoinedAt       time.Time    `gorm:"autoCreateTime" json:"joined_at"`
	LastReadAt     *time.Time   `json:"last_read_at,omitempty"`
	LeftAt         *time.Time   `json:"left_at,omitempty"`
}

