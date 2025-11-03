package models

import (
	"time"
)

type Challenge struct {
	ID          uint `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Sport       string
	Location    string
	CreatorID   uint `gorm:"not null"`
	Creator     User `gorm:"foreignKey:CreatorID"`
	Teams       []Team `gorm:"many2many:challenge_teams;"`
	Users       []User `gorm:"many2many:user_challenges;"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

