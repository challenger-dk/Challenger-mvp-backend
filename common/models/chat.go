package models

import "time"

type Chat struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    // Optional, for named group chats
	Users     []User    `gorm:"many2many:user_chats;"`
	Messages  []Message `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
