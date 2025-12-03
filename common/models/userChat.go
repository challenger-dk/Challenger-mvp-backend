package models

import (
	"time"
)

// UserChat represents the join table between Users and Chats
// It allows us to store metadata like when the user last read the chat.
type UserChat struct {
	ChatID     uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"primaryKey"`
	LastReadAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// TableName overrides the table name to match the existing GORM many2many convention
func (UserChat) TableName() string {
	return "user_chats"
}
