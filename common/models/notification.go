package models

import (
	"time"
)

type NotificationType string

const (
	NotifTypeSystem NotificationType = "system"

	// Team
	NotifTypeTeamInvite      NotificationType = "team_invite"
	NotifTypeTeamAccept      NotificationType = "team_accept"
	NotifTypeTeamDecline     NotificationType = "team_decline"
	NotifTypeTeamRemovedUser NotificationType = "team_removed_user"
	NotifTypeTeamUserLeft    NotificationType = "team_user_left"
	NotifTypeTeamDeleted     NotificationType = "team_deleted"

	// Friend
	NotifTypeFriendReq     NotificationType = "friend_request"
	NotifTypeFriendAccept  NotificationType = "friend_accept"
	NotifTypeFriendDecline NotificationType = "friend_decline"

	// Challenge
	NotifTypeChallengeReq     NotificationType = "challenge_request"
	NotifTypeChallengeAccept  NotificationType = "challenge_accept"
	NotifTypeChallengeDecline NotificationType = "challenge_decline"
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
	InvitationID *uint            // Linked Invitation ID
	IsRead       bool             `gorm:"default:false"`
	IsRelevant   bool             `gorm:"default:true"` // <--- Controls visibility
	CreatedAt    time.Time        `gorm:"autoCreateTime"`
}
