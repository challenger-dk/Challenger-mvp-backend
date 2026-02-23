/*
	Functionality to Create/read/mark/unmark notifications
	Creating notifications does not return any errors, since it should not break the transaction.
	Creating transactions are wrapped in their own transaction, does not break the main transaction.
	If anything goes wrong when creating a notification, it is logged.
*/

package services

import (
	"fmt"
	"log/slog"
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
	ResourceType *models.ResourceType
	InvitationID *uint
}

type NotificationFilters struct {
	IsRead *bool
	Limit  int
	Offset int
}

// Wrapper checking if the notification should be sent before creation.
func CreateNotification(db *gorm.DB, params NotificationParams) {
	// 1. Check if user settings allow this notification
	if !shouldNotify(db, params.RecipientID, params.Type) {
		return
	}

	// 2. Check if Recipient has blocked Actor
	if params.ActorID != nil {
		if IsBlocked(params.RecipientID, *params.ActorID) {
			return
		}
	}

	// Create the actual notification
	persistNotification(db, params)

	// Send push notifications for invitations and selected notification types
	switch params.Type {
	case models.NotifTypeChallengeReq:
		sendInvitationPushNotification(db, params)
	case models.NotifTypeFriendReq:
		sendInvitationPushNotification(db, params)
	case models.NotifTypeFriendAccept:
		sendInvitationPushNotification(db, params)
	case models.NotifTypeChallengeAccept:
		sendInvitationPushNotification(db, params)
	case models.NotifTypeChallengeUpcomming24H, models.NotifTypeChallengeUpcomming1H:
		sendInvitationPushNotification(db, params)
	}
}

// GetMyNotifications fetches notifications with filters.
func GetMyNotifications(userID uint, filters NotificationFilters) ([]models.Notification, error) {
	var notifs []models.Notification

	// Filter: UserID matches AND IsRelevant is true
	// Also exclude notifications from blocked users (actor_id)
	query := config.DB.
		Scopes(ExcludeBlockedUsersOn(userID, "actor_id")).
		Preload("Actor").
		Where("user_id = ? AND is_relevant = ?", userID, true)

	// Apply Read/Unread filter if provided
	if filters.IsRead != nil {
		query = query.Where("is_read = ?", *filters.IsRead)
	}

	// Apply Limit (default to 50 if not specified)
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	} else {
		query = query.Limit(50)
	}

	// Apply Offset
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Order("created_at desc").
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

// HideNotificationByInvitationID marks the notification associated with an invitation as irrelevant.
func HideNotificationByInvitationID(invitationID uint) {
	err := config.DB.Model(&models.Notification{}).
		Where("invitation_id = ?", invitationID).
		Update("is_relevant", false).
		Error

	if err != nil {
		slog.Warn("Failed to hide notification for invitation",
			slog.Uint64("invitation_id", uint64(invitationID)),
			slog.Any("error", err),
		)
	}
}

// ------------- Creators ------------- \\

// ------ INVITATIONS ----- \\

// CreateInvitationNotification sends a notification to the invitee.
// For friend invitations, ResourceID points to the inviter (so the recipient can open the inviter's profile).
func CreateInvitationNotification(db *gorm.DB, inv models.Invitation) {
	switch inv.ResourceType {
	case models.ResourceTypeFriend:
		title := "Ny venneanmodning"
		content := "Du har modtaget en ny venneanmodning."

		// If inviter is preloaded, we can personalize the text
		if inv.Inviter.ID != 0 && inv.Inviter.FirstName != "" {
			content = fmt.Sprintf("%s har inviteret dig til at være venner", inv.Inviter.FirstName)
		}

		rid := inv.InviterId
		rType := models.ResourceTypeFriend

		CreateNotification(db, NotificationParams{
			RecipientID:  inv.InviteeId,
			Type:         models.NotifTypeFriendReq,
			Title:        title,
			Content:      content,
			ActorID:      &inv.InviterId,
			ResourceID:   &rid,
			ResourceType: &rType,
			InvitationID: &inv.ID,
		})
		return

	case models.ResourceTypeTeam:
		resource, err := getResource(inv, db)
		if err != nil {
			slog.Error("Failed to get resource for invitation notification",
				slog.Int("invitation_id", int(inv.ID)),
				slog.Any("error", err),
			)
			return
		}

		team, ok := resource.(models.Team)
		if !ok {
			slog.Error("Failed to cast resource to team",
				slog.Int("invitation_id", int(inv.ID)),
			)
			return
		}

		title := "Ny klubinvitation"
		content := fmt.Sprintf("Du er blevet inviteret til at blive en del af: %s", team.Name)
		rid := inv.ResourceID
		rType := inv.ResourceType

		CreateNotification(db, NotificationParams{
			RecipientID:  inv.InviteeId,
			Type:         models.NotifTypeTeamInvite,
			Title:        title,
			Content:      content,
			ActorID:      &inv.InviterId,
			ResourceID:   &rid,
			ResourceType: &rType,
			InvitationID: &inv.ID,
		})
		return

	case models.ResourceTypeChallenge:
		// If you later want the challenge name in the notification, you can use getResource here.
		title := "Du er inviteret til en udfordring"
		content := "Du er blevet inviteret – vil du være med?"
		rid := inv.ResourceID
		rType := inv.ResourceType

		CreateNotification(db, NotificationParams{
			RecipientID:  inv.InviteeId,
			Type:         models.NotifTypeChallengeReq,
			Title:        title,
			Content:      content,
			ActorID:      &inv.InviterId,
			ResourceID:   &rid,
			ResourceType: &rType,
			InvitationID: &inv.ID,
		})
		return

	default:
		slog.Warn("Notification skipped: unknown resource type",
			slog.String("resource_type", string(inv.ResourceType)),
		)
		return
	}
}

// CreateAcceptedInvitationNotification sends notification to the inviter.
// For friend accept, ResourceID points to the invitee (the person who accepted), so inviter can open their profile.
func CreateAcceptedInvitationNotification(db *gorm.DB, inv models.Invitation) {
	switch inv.ResourceType {
	case models.ResourceTypeTeam:
		title := "Klubinvitation accepteret"
		content := fmt.Sprintf("%s er blevet medlem af din klub", inv.Invitee.FirstName)

		rid := inv.ResourceID
		rType := inv.ResourceType

		CreateNotification(db, NotificationParams{
			RecipientID:  inv.InviterId,
			Type:         models.NotifTypeTeamAccept,
			Title:        title,
			Content:      content,
			ActorID:      &inv.InviteeId,
			ResourceID:   &rid,
			ResourceType: &rType,
			InvitationID: &inv.ID,
		})
		return

	case models.ResourceTypeFriend:
		title := "Venneanmodning accepteret"
		content := fmt.Sprintf("%s har accepteret din venneanmodning", inv.Invitee.FirstName)

		rid := inv.InviteeId
		rType := models.ResourceTypeFriend

		CreateNotification(db, NotificationParams{
			RecipientID:  inv.InviterId,
			Type:         models.NotifTypeFriendAccept,
			Title:        title,
			Content:      content,
			ActorID:      &inv.InviteeId,
			ResourceID:   &rid,
			ResourceType: &rType,
			InvitationID: &inv.ID,
		})
		return

	case models.ResourceTypeChallenge:
		title := "Udfordringsinvitation accepteret"
		content := fmt.Sprintf("%s har accepteret din challenge invitation", inv.Invitee.FirstName)

		rid := inv.ResourceID
		rType := inv.ResourceType

		CreateNotification(db, NotificationParams{
			RecipientID:  inv.InviterId,
			Type:         models.NotifTypeChallengeAccept,
			Title:        title,
			Content:      content,
			ActorID:      &inv.InviteeId,
			ResourceID:   &rid,
			ResourceType: &rType,
			InvitationID: &inv.ID,
		})
		return

	default:
		slog.Warn("Notification skipped: unknown resource type",
			slog.String("resource_type", string(inv.ResourceType)),
		)
		return
	}
}

// CreateDeclinedInvitationNotification sends notification to the inviter.
// For friend decline, ResourceID points to the invitee (the person who declined), so inviter can open their profile.
func CreateDeclinedInvitationNotification(db *gorm.DB, inv models.Invitation) {
	switch inv.ResourceType {
	case models.ResourceTypeTeam:
		title := "Klubinvitation afvist"
		content := fmt.Sprintf("%s har afvist din klubinvitation", inv.Invitee.FirstName)

		rid := inv.ResourceID
		rType := inv.ResourceType

		CreateNotification(db, NotificationParams{
			RecipientID:  inv.InviterId,
			Type:         models.NotifTypeTeamDecline,
			Title:        title,
			Content:      content,
			ActorID:      &inv.InviteeId,
			ResourceID:   &rid,
			ResourceType: &rType,
			InvitationID: &inv.ID,
		})
		return

	case models.ResourceTypeFriend:
		title := "Venneanmodning afvist"
		content := fmt.Sprintf("%s har afvist din venneanmodning", inv.Invitee.FirstName)

		rid := inv.InviteeId
		rType := models.ResourceTypeFriend

		CreateNotification(db, NotificationParams{
			RecipientID:  inv.InviterId,
			Type:         models.NotifTypeFriendDecline,
			Title:        title,
			Content:      content,
			ActorID:      &inv.InviteeId,
			ResourceID:   &rid,
			ResourceType: &rType,
			InvitationID: &inv.ID,
		})
		return

	case models.ResourceTypeChallenge:
		title := "Udfordringsinvitation afvist"
		content := fmt.Sprintf("%s har afvist din challenge invitation", inv.Invitee.FirstName)

		rid := inv.ResourceID
		rType := inv.ResourceType

		CreateNotification(db, NotificationParams{
			RecipientID:  inv.InviterId,
			Type:         models.NotifTypeChallengeDecline,
			Title:        title,
			Content:      content,
			ActorID:      &inv.InviteeId,
			ResourceID:   &rid,
			ResourceType: &rType,
			InvitationID: &inv.ID,
		})
		return

	default:
		slog.Warn("Notification skipped: unknown resource type",
			slog.String("resource_type", string(inv.ResourceType)),
		)
		return
	}
}

// ------ TEAMS ----- \\

// User removed from team
func CreateRemovedUserFromTeamNotification(db *gorm.DB, userID uint, team models.Team) {
	title := "Du er blevet fjernet fra en klub"
	content := fmt.Sprintf("Du er blevet fjernet fra '%s'", team.Name)

	rid := team.ID
	rType := models.ResourceTypeTeam

	CreateNotification(db, NotificationParams{
		RecipientID:  userID,
		Type:         models.NotifTypeTeamRemovedUser,
		Title:        title,
		Content:      content,
		ResourceID:   &rid,
		ResourceType: &rType,
	})
}

// User left team, notifies the creator only
func CreateUserLeftTeamNotification(db *gorm.DB, leaver models.User, team models.Team) {
	title := "En bruger har forladt klubben"
	content := fmt.Sprintf("%s har forladt '%s'", leaver.FirstName, team.Name)

	rid := team.ID
	rType := models.ResourceTypeTeam

	CreateNotification(db, NotificationParams{
		RecipientID:  team.CreatorID,
		Type:         models.NotifTypeTeamUserLeft,
		Title:        title,
		Content:      content,
		ActorID:      &leaver.ID,
		ResourceID:   &rid,
		ResourceType: &rType,
	})
}

func CreateTeamDeletedNotification(db *gorm.DB, user models.User, team models.Team) {
	title := "En klub du er medlem af er blevet slettet"
	content := fmt.Sprintf("'%s' er blevet slettet", team.Name)

	rid := team.ID
	rType := models.ResourceTypeTeam

	CreateNotification(db, NotificationParams{
		RecipientID:  user.ID,
		Type:         models.NotifTypeTeamDeleted,
		Title:        title,
		Content:      content,
		ResourceID:   &rid,
		ResourceType: &rType,
	})
}

// ------ CHALLENGES ----- \\

func CreateUserJoinedChallengeNotificationToCreator(db *gorm.DB, user models.User, challenge models.Challenge) {
	title := "Ny deltager har tilmeldt sig din udfordring"
	content := fmt.Sprintf("%s har tilmeldt sig '%s'", user.FirstName, challenge.Name)

	rid := challenge.ID
	rType := models.ResourceTypeChallenge

	CreateNotification(db, NotificationParams{
		RecipientID:  challenge.CreatorID,
		Type:         models.NotifTypeChallengeJoin,
		Title:        title,
		Content:      content,
		ResourceID:   &rid,
		ResourceType: &rType,
	})
}

func CreateUserJoinedChallengeNotification(db *gorm.DB, user models.User, challenge models.Challenge) {
	title := "Du har tilmeldt dig en udfordring"
	content := fmt.Sprintf("Du har tilmeldt dig '%s'", challenge.Name)

	rid := challenge.ID
	rType := models.ResourceTypeChallenge

	CreateNotification(db, NotificationParams{
		RecipientID:  user.ID,
		Type:         models.NotifTypeChallengeJoin,
		Title:        title,
		Content:      content,
		ResourceID:   &rid,
		ResourceType: &rType,
	})
}

func CreateUserLeftChallengeNotification(db *gorm.DB, user models.User, challenge models.Challenge) {
	title := "En deltager har forladt udfordringen"
	content := fmt.Sprintf("%s har forladt '%s'", user.FirstName, challenge.Name)

	rid := challenge.ID
	rType := models.ResourceTypeChallenge

	CreateNotification(db, NotificationParams{
		RecipientID:  challenge.CreatorID,
		Type:         models.NotifTypeChallengeUserLeft,
		Title:        title,
		Content:      content,
		ResourceID:   &rid,
		ResourceType: &rType,
	})
}

func CreateNotificationUpcomingChallenge(db *gorm.DB, user models.User, challenge models.Challenge, notifType models.NotificationType) {
	title := "Din udfordring starter snart!"
	content := fmt.Sprintf("Din udfordring starter kl. '%s'", challenge.StartTime.Format("15:04")) // default content
	// Set content based on notification type
	switch notifType {
	case models.NotifTypeChallengeUpcomming24H:
		content = fmt.Sprintf("Din udfordring starter om 24 timer kl. '%s'", challenge.StartTime.Format("15:04"))
	case models.NotifTypeChallengeUpcomming1H:
		content = fmt.Sprintf("Din udfordring starter om 1 time kl. '%s'", challenge.StartTime.Format("15:04"))
	default:
		slog.Warn("CreateNotificationUpcomingChallenge called with invalid notifType",
			slog.String("notif_type", string(notifType)),
		)
		return
	}

	rid := challenge.ID
	rType := models.ResourceTypeChallenge

	CreateNotification(db, NotificationParams{
		RecipientID:  user.ID,
		Type:         notifType,
		Title:        title,
		Content:      content,
		ResourceID:   &rid,
		ResourceType: &rType,
	})
}

func CreateNotificationChallengeFullParticipation(db *gorm.DB, user models.User, challenge models.Challenge) {
	title := "Din udfordring har fuld deltagelse!"
	content := "Alle deltagere er nu klar"

	rid := challenge.ID
	rType := models.ResourceTypeChallenge

	CreateNotification(db, NotificationParams{
		RecipientID:  user.ID,
		Type:         models.NotifTypeChallengeFullParticipation,
		Title:        title,
		Content:      content,
		ResourceID:   &rid,
		ResourceType: &rType,
	})
}

func CreateNotificationChallengeNotAnswered24H(db *gorm.DB, user models.User, challenge models.Challenge) {
	title := "Husk at svare på invitation"
	content := "Du er inviteret til en udfordring – husk at svare inden start."

	rid := challenge.ID
	rType := models.ResourceTypeChallenge

	CreateNotification(db, NotificationParams{
		RecipientID:  user.ID,
		Type:         models.NotifTypeChallengeNotAnswered24H,
		Title:        title,
		Content:      content,
		ResourceID:   &rid,
		ResourceType: &rType,
	})
}

func CreateNotificationChallengeMissingParticipants(db *gorm.DB, user models.User, challenge models.Challenge) {
	title := "Din udfordring mangler deltagere"
	content := "Din udfordring mangler deltagere – husk at invitere flere."

	rid := challenge.ID
	rType := models.ResourceTypeChallenge

	CreateNotification(db, NotificationParams{
		RecipientID:  user.ID,
		Type:         models.NotifTypeChallengeMissingParticipants,
		Title:        title,
		Content:      content,
		ResourceID:   &rid,
		ResourceType: &rType,
	})
}

// -------------- Private -------------- \\
func shouldNotify(db *gorm.DB, userID uint, notifType models.NotificationType) bool {
	var settings models.UserSettings

	// If settings row doesn't exist (or any error), default to sending notifications.
	if err := db.First(&settings, userID).Error; err != nil {
		return true
	}

	// Look up a setting getter for this notification type.
	if getter, ok := notificationSettingsMap[notifType]; ok {
		return getter(settings)
	}

	// Unknown/unmapped types default to "send".
	return true
}

// sendInvitationPushNotification sends a push notification to the invitee for friend requests and challenge invitations.
// Uses db (may be a transaction) to load the recipient so it works correctly when called from within a transaction.
// Errors are logged but do not affect the caller.
func sendInvitationPushNotification(db *gorm.DB, params NotificationParams) {
	var recipient models.User
	if err := db.Select("id", "expo_token").
		First(&recipient, params.RecipientID).Error; err != nil {
		slog.Warn("Failed to load recipient for invitation push",
			slog.Uint64("recipient_id", uint64(params.RecipientID)),
			slog.String("notification_type", string(params.Type)),
			slog.Any("error", err),
		)
		return
	}
	if recipient.ExpoToken == "" {
		slog.Info("Skipping invitation push: recipient has no expo token",
			slog.Uint64("recipient_id", uint64(params.RecipientID)),
			slog.String("notification_type", string(params.Type)),
		)
		return
	}

	data := map[string]any{}
	if params.ResourceType != nil {
		data["resource_type"] = string(*params.ResourceType)
	}
	if params.ResourceID != nil {
		data["resource_id"] = *params.ResourceID
	}
	if params.InvitationID != nil {
		data["invitation_id"] = *params.InvitationID
	}

	err := SendExpoPushNotification(recipient.ExpoToken, params.Title, params.Content, data)
	if err != nil {
		slog.Warn("Failed to send invitation push notification",
			slog.Uint64("recipient_id", uint64(params.RecipientID)),
			slog.String("notification_type", string(params.Type)),
			slog.Any("error", err),
		)
		if IsDeviceNotRegistered(err) {
			config.DB.Model(&models.User{}).Where("id = ?", recipient.ID).Update("expo_token", "")
		}
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
		InvitationID: params.InvitationID,
		IsRelevant:   true,
	}

	// Savepoint
	err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&n).Error
	})

	if err != nil {
		slog.Error("Notification DB error (swallowed)", slog.Any("error", err))
	}
}

// notificationSettingsMap maps each notification type to the relevant user-setting toggle.
var notificationSettingsMap = map[models.NotificationType]func(s models.UserSettings) bool{
	// System
	models.NotifTypeSystem: func(s models.UserSettings) bool { return true },

	// ---------------- Team ----------------
	models.NotifTypeTeamInvite: func(s models.UserSettings) bool { return s.NotifyTeamInvites },

	models.NotifTypeTeamAccept:      func(s models.UserSettings) bool { return s.NotifyTeamMembership },
	models.NotifTypeTeamDecline:     func(s models.UserSettings) bool { return s.NotifyTeamMembership },
	models.NotifTypeTeamRemovedUser: func(s models.UserSettings) bool { return s.NotifyTeamMembership },
	models.NotifTypeTeamUserLeft:    func(s models.UserSettings) bool { return s.NotifyTeamMembership },
	models.NotifTypeTeamDeleted:     func(s models.UserSettings) bool { return s.NotifyTeamMembership },

	// ---------------- Friend ----------------
	models.NotifTypeFriendReq:     func(s models.UserSettings) bool { return s.NotifyFriendRequests },
	models.NotifTypeFriendAccept:  func(s models.UserSettings) bool { return s.NotifyFriendUpdates },
	models.NotifTypeFriendDecline: func(s models.UserSettings) bool { return s.NotifyFriendUpdates },

	// ---------------- Challenge ----------------
	models.NotifTypeChallengeReq:     func(s models.UserSettings) bool { return s.NotifyChallengeInvites },
	models.NotifTypeChallengeAccept:  func(s models.UserSettings) bool { return s.NotifyChallengeInvites },
	models.NotifTypeChallengeDecline: func(s models.UserSettings) bool { return s.NotifyChallengeInvites },

	models.NotifTypeChallengeCreated:             func(s models.UserSettings) bool { return s.NotifyChallengeUpdates },
	models.NotifTypeChallengeJoin:                func(s models.UserSettings) bool { return s.NotifyChallengeUpdates },
	models.NotifTypeChallengeUserLeft:            func(s models.UserSettings) bool { return s.NotifyChallengeUpdates },
	models.NotifTypeChallengeFullParticipation:   func(s models.UserSettings) bool { return s.NotifyChallengeUpdates },
	models.NotifTypeChallengeMissingParticipants: func(s models.UserSettings) bool { return s.NotifyChallengeUpdates },

	models.NotifTypeChallengeUpcomming24H:   func(s models.UserSettings) bool { return s.NotifyChallengeReminders },
	models.NotifTypeChallengeUpcomming1H:    func(s models.UserSettings) bool { return s.NotifyChallengeReminders },
	models.NotifTypeChallengeNotAnswered24H: func(s models.UserSettings) bool { return s.NotifyChallengeReminders },
}
