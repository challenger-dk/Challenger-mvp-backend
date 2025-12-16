package models

import (
	"time"
)

type ConversationType string

const (
	ConversationTypeDirect ConversationType = "direct"
	ConversationTypeGroup  ConversationType = "group"
	ConversationTypeTeam   ConversationType = "team"
)

type Conversation struct {
	ID        uint             `gorm:"primaryKey" json:"id"`
	Type      ConversationType `gorm:"type:VARCHAR(20);not null;check:type IN ('direct','group','team')" json:"type"`
	Title     *string          `gorm:"type:VARCHAR(255)" json:"title,omitempty"` // For group chats
	TeamID    *uint            `gorm:"uniqueIndex:idx_team_conversation" json:"team_id,omitempty"`
	Team      *Team            `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	DirectKey *string          `gorm:"type:VARCHAR(100);uniqueIndex:idx_direct_key" json:"-"` // For direct chat uniqueness
	CreatedAt time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time        `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Participants []ConversationParticipant `gorm:"foreignKey:ConversationID" json:"participants,omitempty"`
	Messages     []Message                 `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}
