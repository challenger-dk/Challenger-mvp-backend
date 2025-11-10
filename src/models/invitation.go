package models

import (
	"time"
)

type Invitation struct {
	ID           uint `gorm:"primaryKey"`
	InviterId    uint `gorm:"not null"`
	InviteeId    uint `gorm:"not null"`
	Note         string
	ResourceType InvitationType   `gorm:"type:VARCHAR(20);not null;check:resource_type IN ('team')"`
	ResourceID   uint             `gorm:"not null"`
	Status       InvitationStatus `gorm:"type:VARCHAR(20);not null;check:status IN ('pending','accepted','declined')"`
	CreatedAt    time.Time        `gorm:"autoCreateTime"`
	UpdatedAt    time.Time        `gorm:"autoUpdateTime"`
}

type InvitationType string
type InvitationStatus string

// Only allowed resource types and statuses
const (
	ResourceTypeTeam InvitationType = "team"
)

const (
	StatusPending  InvitationStatus = "pending"
	StatusAccepted InvitationStatus = "accepted"
	StatusDeclined InvitationStatus = "declined"
)
