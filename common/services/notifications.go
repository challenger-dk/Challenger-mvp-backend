/*
	Functionality to Create/read/mark/unmark notifications
	Creating notifications does not return any errors, since it should not break the transaction.
	Creating transactions are wrapped in their own transaction, does not break the main transaction.
	If anything goes wrong when creating a notification, it is logged.
*/

package services

import (
	"fmt"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

type NotificationParams struct {
	RecipientID  uint
	Type         models.NotificationType
	Title        string
	Content      string
	ActorID      *uint
	ResourceID   *uint
	ResourceType *string
}

// Wrapper checking if the notification should be sent before creation.
func CreateNotification(db *gorm.DB, params NotificationParams) {
	if !shouldNotify(db, params.RecipientID, params.Type) {
		return
	}

	// Create the actual notification
	persistNotification(db, params)
}

// GetMyNotifications fetches notifications.
func GetMyNotifications(userID uint) ([]models.Notification, error) {
	var notifs []models.Notification
	err := config.DB.Preload("Actor").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(50).
		Find(&notifs).
		Error

	return notifs, err
}

// MarkNotificationAsRead marks a single notification as read.
func MarkNotificationAsRead(notifID, userID uint) error {
	return config.DB.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("is_read", true).
		Error
}

// MarkAllNotificationsAsRead marks all user's notifications as read.
func MarkAllNotificationsAsRead(userID uint) error {
	return config.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).
		Error
}

// ------------- Creators ------------- \\

// Sends notification to the invitee
func CreateInvitationNotification(db *gorm.DB, inv models.Invitation) {
	var title, content string
	var notifType models.NotificationType

	switch inv.ResourceType {
	case models.ResourceTypeTeam:
		title = "Team Invitation"
		content = "You have been invited to join a team."
		notifType = models.NotifTypeTeamInvite

	case models.ResourceTypeFriend:
		title = "Friend Request"
		content = "You have a new friend request."
		notifType = models.NotifTypeFriendReq

	case models.ResourceTypeChallenge:
		title = "Challenge Invitation"
		content = "You have been invited to join a challenge."
		notifType = models.NotifTypeChallengeReq

	default:
		fmt.Printf("⚠️ Notification skipped: Unknown resource type %s\n", inv.ResourceType)
		return
	}

	rType := string(inv.ResourceType)

	CreateNotification(db, NotificationParams{
		RecipientID:  inv.InviteeId,
		Type:         notifType,
		Title:        title,
		Content:      content,
		ActorID:      &inv.InviterId,
		ResourceID:   &inv.ResourceID,
		ResourceType: &rType,
	})
}

// Sends notification to the one who sent the invitation
func CreateAcceptedInvitationNotification(db *gorm.DB, inv models.Invitation) {
	var title, content string
	var notifType models.NotificationType

	switch inv.ResourceType {
	case models.ResourceTypeTeam:
		title = "Team Invitation Accepted"
		content = fmt.Sprintf("%s has joined your team", inv.Invitee.FirstName)
		notifType = models.NotifTypeTeamAccept

	case models.ResourceTypeFriend:
		title = "Friend Request Accepted"
		content = fmt.Sprintf("%s has accepted your friend request", inv.Invitee.FirstName)
		notifType = models.NotifTypeFriendAccept

	case models.ResourceTypeChallenge:
		title = "Challenge Invitation Accepted"
		content = fmt.Sprintf("%s has accepted your challenge invitation", inv.Invitee.FirstName)
		notifType = models.NotifTypeChallengeAccept

	default:
		fmt.Printf("⚠️ Notification skipped: Unknown resource type %s\n", inv.ResourceType)
		return
	}

	rType := string(inv.ResourceType)

	CreateNotification(db, NotificationParams{
		RecipientID:  inv.InviterId,
		Type:         notifType,
		Title:        title,
		Content:      content,
		ActorID:      &inv.InviteeId,
		ResourceID:   &inv.ResourceID,
		ResourceType: &rType,
	})
}

// Sends notification to the one who sent the invitation
func CreateDeclinedInvitationNotification(db *gorm.DB, inv models.Invitation) {
	var title, content string
	var notifType models.NotificationType

	switch inv.ResourceType {
	case models.ResourceTypeTeam:
		title = "Team Invitation Declined"
		content = fmt.Sprintf("%s has declined your team invitation", inv.Invitee.FirstName)
		notifType = models.NotifTypeTeamDecline

	case models.ResourceTypeFriend:
		title = "Friend Request Declined"
		content = fmt.Sprintf("%s has declined your friend request", inv.Invitee.FirstName)
		notifType = models.NotifTypeFriendDecline

	case models.ResourceTypeChallenge:
		title = "Challenge Invitation Declined"
		content = fmt.Sprintf("%s has declined your challenge invitation", inv.Invitee.FirstName)
		notifType = models.NotifTypeChallengeDecline

	default:
		fmt.Printf("⚠️ Notification skipped: Unknown resource type %s\n", inv.ResourceType)
		return
	}

	rType := string(inv.ResourceType)

	CreateNotification(db, NotificationParams{
		RecipientID:  inv.InviterId,
		Type:         notifType,
		Title:        title,
		Content:      content,
		ActorID:      &inv.InviteeId,
		ResourceID:   &inv.ResourceID,
		ResourceType: &rType,
	})
}

// -------------- Private -------------- \\
func shouldNotify(db *gorm.DB, userID uint, notifType models.NotificationType) bool {
	var settings models.UserSettings

	if err := db.First(&settings, userID).Error; err != nil {
		return true
	}

	switch notifType {
	// Team
	case models.NotifTypeTeamInvite:
		return settings.NotifyTeamInvite
	case models.NotifTypeTeamAccept:
		return settings.NotifyTeamInvite

	// Friend
	case models.NotifTypeFriendReq:
		return settings.NotifyFriendReq
	case models.NotifTypeFriendAccept:
		return settings.NotifyFriendReq

	default:
		return true
	}
}

// persistNotification performs the actual DB write.
// It is private so no one accidentally bypasses the checks.
func persistNotification(db *gorm.DB, params NotificationParams) {
	n := models.Notification{
		UserID:       params.RecipientID,
		Type:         params.Type,
		Title:        params.Title,
		Content:      params.Content,
		ActorID:      params.ActorID,
		ResourceID:   params.ResourceID,
		ResourceType: params.ResourceType,
	}

	// Savepoint
	err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&n).Error
	})

	if err != nil {
		fmt.Printf("⚠️ Notification DB Error (swallowed): %v\n", err)
	}
}
