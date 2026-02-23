package models

import (
	"server/common/models/types"
)

// Facility represents a sports facility (e.g. courts, halls) that can host challenges.
// Data is imported from facilities.json. Facilities can have many challenges.
type Facility struct {
	ID           uint        `gorm:"primaryKey"`
	ExternalID   string      `gorm:"type:text;uniqueIndex;not null" json:"id"` // Original ID from facilities.json (e.g. "facility-aabenraa-1")
	Name         string      `gorm:"type:text;not null"`
	DetailedName string      `gorm:"type:text;not null"`
	Address      string      `gorm:"type:text;not null"`
	Website      *string     `gorm:"type:text;default:null"`
	Email        *string     `gorm:"type:text;default:null"`
	FacilityType string      `gorm:"type:text;not null"`
	Indoor       bool        `gorm:"default:false"`
	Notes        *string     `gorm:"type:text;default:null"`
	Coordinates  types.Point `gorm:"type:geography(Point,4326);not null"`
	PostalCode   string      `gorm:"type:text;not null"`
	City         string      `gorm:"type:text;not null"`
	Country      string      `gorm:"type:text;not null"`
	Challenges   []Challenge `gorm:"foreignKey:FacilityID"`
}

// TableName overrides the table name for Facility
func (Facility) TableName() string {
	return "facilities"
}
