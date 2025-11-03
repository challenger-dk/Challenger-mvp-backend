package models

import (
	"time"
)

type Team struct {
	ID        uint `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Users     []User `gorm:"many2many:user_teams;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}