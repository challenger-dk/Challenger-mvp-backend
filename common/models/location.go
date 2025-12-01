package models

import (
	"time"

	"server/common/models/types"
)

// Publicly exposed, if needed use outside of package
type Point = types.Point

type Location struct {
	ID          uint      `gorm:"primaryKey"`
	Address     string    `gorm:"not null"`
	Coordinates Point     `gorm:"type:geography(Point,4326);not null;uniqueIndex"`
	PostalCode  string    `gorm:"not null"`
	City        string    `gorm:"not null"`
	Country     string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
