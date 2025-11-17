package models

import (
	"server/models/types"
	"time"
)

type Location struct {
	ID          uint        `gorm:"primaryKey"`
	Address     string      `gorm:"not null"`
	Coordinates types.Point `gorm:"type:geography(Point,4326);not null;uniqueIndex"`
	PostalCode  string      `gorm:"not null"`
	City        string      `gorm:"not null"`
	Country     string      `gorm:"not null"`
	CreatedAt   time.Time   `gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime"`
}
