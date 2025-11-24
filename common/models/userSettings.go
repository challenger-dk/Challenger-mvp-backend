package models

import "time"

type UserSettings struct {
	UserID uint `gorm:"primaryKey"`

	// ------ Notification Settings ------ \\
	// Team
	NotifyTeamInvite bool `gorm:"default:true"`

	// Friends
	NotifyFriendReq bool `gorm:"default:true"`

	// Challenges
	NotifyChallengeInvite bool `gorm:"default:true"`
	NotifyChallengeUpdate bool `gorm:"default:true"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
