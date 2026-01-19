package models

import "time"

type UserSettings struct {
	UserID uint `gorm:"primaryKey"`

	// --------- Team notifications --------- \\

	// Affects:
	// - team_invite
	NotifyTeamInvites bool `gorm:"default:true"`

	// Affects:
	// - team_accept
	// - team_decline
	// - team_removed_user
	// - team_user_left
	// - team_deleted
	NotifyTeamMembership bool `gorm:"default:true"`

	// --------- Friend notifications --------- \\

	// Affects:
	// - friend_request
	NotifyFriendRequests bool `gorm:"default:true"`

	// Affects:
	// - friend_accept
	// - friend_decline
	NotifyFriendUpdates bool `gorm:"default:true"`

	// --------- Challenge notifications --------- \\

	// Affects:
	// - challenge_request
	// - challenge_accept
	// - challenge_decline
	NotifyChallengeInvites bool `gorm:"default:true"`

	// Affects:
	// - challenge_created
	// - challenge_joined
	// - challenge_user_left
	// - challenge_full_participation
	// - challenge_missing_participants
	NotifyChallengeUpdates bool `gorm:"default:true"`

	// Affects:
	// - challenge_upcomming_24h
	// - challenge_upcomming_1h
	// - challenge_invitation_not_answered_24h
	NotifyChallengeReminders bool `gorm:"default:true"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
