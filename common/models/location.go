package models

import (
	"time"

	"server/common/models/types"
)

// Publicly exposed, if needed use outside of package
type Point = types.Point

type LocationType string

const (
	LocationTypeAddress  LocationType = "address"
	LocationTypeFacility LocationType = "facility"
)

type Location struct {
	ID          uint         `gorm:"primaryKey"`
	Type        LocationType `gorm:"type:VARCHAR(20);not null;default:'address';check:type IN ('address','facility')"`
	Address     string       `gorm:"not null"`
	Coordinates Point        `gorm:"type:geography(Point,4326);not null;uniqueIndex"`
	PostalCode  string       `gorm:"not null"`
	City        string       `gorm:"not null"`
	Country     string       `gorm:"not null"`
	// Facility-specific fields (nullable)
	FacilityID   *string   `gorm:"default:null"`
	FacilityName *string   `gorm:"default:null"`
	DetailedName *string   `gorm:"default:null"`
	Email        *string   `gorm:"default:null"`
	Website      *string   `gorm:"default:null"`
	FacilityType *string   `gorm:"default:null"`
	Indoor       *bool     `gorm:"default:null"`
	Notes        *string   `gorm:"default:null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}
