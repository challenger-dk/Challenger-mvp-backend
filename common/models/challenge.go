package models

import (
	"time"

	"gorm.io/gorm"
)

type Challenge struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Sport       string
	Location    Location       `gorm:"foreignKey:LocationID"`
	LocationID  uint           `gorm:"not null"`
	CreatorID   uint           `gorm:"not null"`
	Creator     User           `gorm:"foreignKey:CreatorID"`
	Teams       []Team         `gorm:"many2many:challenge_teams;"`
	Users       []User         `gorm:"many2many:user_challenges;"`
	IsIndoor    bool           `gorm:"default:false"`
	IsPublic    bool           `gorm:"default:false"`
	IsCompleted bool           `gorm:"default:false"`
	PlayFor     *string        `gorm:"default:null"`
	HasCost     bool           `gorm:"default:false"`
	Comment     *string        `gorm:"default:null"`
	TeamSize    *int           `gorm:"default:null"`
	Date        time.Time      `gorm:"not null"`
	StartTime   time.Time      `gorm:"not null"`
	EndTime     *time.Time     `gorm:"default:null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
