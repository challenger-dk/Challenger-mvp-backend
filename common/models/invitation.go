package models

import (
	"time"
)

type InvitationStatus string

const (
	StatusPending  InvitationStatus = "pending"
	StatusAccepted InvitationStatus = "accepted"
	StatusDeclined InvitationStatus = "declined"
)

type Invitation struct {
	ID           uint `gorm:"primaryKey"`
	InviterId    uint `gorm:"not null;uniqueIndex:idx_unique_invitation"`
	Inviter      User `gorm:"foreignKey:InviterId"`
	InviteeId    uint `gorm:"not null;uniqueIndex:idx_unique_invitation"`
	Invitee      User `gorm:"foreignKey:InviteeId"`
	Note         string
	ResourceType ResourceType     `gorm:"type:VARCHAR(20);not null;check:resource_type IN ('team','friend','challenge');uniqueIndex:idx_unique_invitation"`
	ResourceID   uint             `gorm:"not null;uniqueIndex:idx_unique_invitation"`
	Status       InvitationStatus `gorm:"type:VARCHAR(20);not null;default:pending;check:status IN ('pending','accepted','declined')"` // Defualt 'pending'
	CreatedAt    time.Time        `gorm:"autoCreateTime"`
	UpdatedAt    time.Time        `gorm:"autoUpdateTime"`
}
