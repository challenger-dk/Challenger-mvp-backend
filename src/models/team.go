package models

import (
	"time"
)

type Team struct {
	ID uint `gorm:"primaryKey"`

	// Attributes
	Name        string  `gorm:"not null"`
	Description *string `` // nullable

	// Relationships
	Sports     []Sport   `gorm:"many2many:team_sports;"`
	Users      []User    `gorm:"many2many:user_teams;"`
	Creator    User      `gorm:"foreignKey:CreatorID"`
	CreatorID  uint      `gorm:"not null"`
	LocationID *uint     `gorm:"index"` // Nullable foreign key
	Location   *Location `gorm:"foreignKey:LocationID"`

	//Other
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
