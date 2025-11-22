package models

import (
	"time"
)

type NotificationType string

const (
	NotifTypeSystem NotificationType = "system"

	// Team
	NotifTypeTeamInvite  NotificationType = "team_invite"
	NotifTypeTeamAccept  NotificationType = "team_accept"
	NotifTypeTeamDecline NotificationType = "team_decline"

	// Friend
	NotifTypeFriendReq     NotificationType = "friend_request"
	NotifTypeFriendAccept  NotificationType = "friend_accept"
	NotifTypeFriendDecline NotificationType = "friend_decline"

	// Challenge
	NotifTypeChallengeReq NotificationType = "challenge_request"
)

type Notification struct {
	ID           uint             `gorm:"primaryKey"`
	UserID       uint             `gorm:"not null;index"` // The recipient
	User         User             `gorm:"foreignKey:UserID"`
	ActorID      *uint            `gorm:"index"` // Who triggered it (optional)
	Actor        *User            `gorm:"foreignKey:ActorID"`
	Type         NotificationType `gorm:"type:VARCHAR(50);not null"`
	Title        string           `gorm:"not null"`
	Content      string           `gorm:"not null"`
	ResourceID   *uint            // Linked ID (e.g., TeamID)
	ResourceType *string          // Linked Type (e.g., "team")
	IsRead       bool             `gorm:"default:false"`
	CreatedAt    time.Time        `gorm:"autoCreateTime"`
}
