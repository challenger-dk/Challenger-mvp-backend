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
	ResourceType *string
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
}

// GetMyNotifications fetches notifications with filters.
func GetMyNotifications(userID uint, filters NotificationFilters) ([]models.Notification, error) {
	var notifs []models.Notification

	// Filter: UserID matches AND IsRelevant is true
	query := config.DB.Preload("Actor").
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
			slog.Int("invitation_id", int(invitationID)),
			slog.Any("error", err),
		)
	}
}

// ------------- Creators ------------- \\

// ------ INVITATIONS ----- \\
// Sends notification to the invitee
func CreateInvitationNotification(db *gorm.DB, inv models.Invitation) {
	var title, content string
	var notifType models.NotificationType

	var resource, err = getResource(inv, db)
	if err != nil {
		slog.Error("Failed to get resource for invitation notification",
			slog.Int("invitation_id", int(inv.ID)),
			slog.Any("error", err),
		)
		return
	}

	var resourceName string
	if team, ok := resource.(models.Team); ok {
		resourceName = team.Name
	}

	switch inv.ResourceType {
	case models.ResourceTypeTeam:
		title = "Ny klubinvitation"
		content = fmt.Sprintf("Du er blevet inviteret til at blive en del af: %s", resourceName)
		notifType = models.NotifTypeTeamInvite

	case models.ResourceTypeFriend:
		title = "Ny venneanmodning"
		content = "Du har modtaget en ny venneanmodning."
		notifType = models.NotifTypeFriendReq

	case models.ResourceTypeChallenge:
		title = "Du er inviteret til en udfordring"
		content = "Du er blevet inviteret – vil du være med?"
		notifType = models.NotifTypeChallengeReq

	default:
		slog.Warn("Notification skipped: unknown resource type", slog.String("resource_type", string(inv.ResourceType)))
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
		InvitationID: &inv.ID,
	})
}

// Sends notification to the one who sent the invitation
func CreateAcceptedInvitationNotification(db *gorm.DB, inv models.Invitation) {
	var title, content string
	var notifType models.NotificationType

	switch inv.ResourceType {
	case models.ResourceTypeTeam:
		title = "Holdinvitation accepteret"
		content = fmt.Sprintf("%s er blevet medlem af dit hold", inv.Invitee.FirstName)
		notifType = models.NotifTypeTeamAccept

	case models.ResourceTypeFriend:
		title = "Venneanmodning accepteret"
		content = fmt.Sprintf("%s har accepteret din venneanmodning", inv.Invitee.FirstName)
		notifType = models.NotifTypeFriendAccept

	case models.ResourceTypeChallenge:
		title = "Udfordringsinvitation accepteret"
		content = fmt.Sprintf("%s har accepteret din challenge invitation", inv.Invitee.FirstName)
		notifType = models.NotifTypeChallengeAccept

	default:
		slog.Warn("Notification skipped: unknown resource type", slog.String("resource_type", string(inv.ResourceType)))
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
		InvitationID: &inv.ID,
	})
}

// Sends notification to the one who sent the invitation
func CreateDeclinedInvitationNotification(db *gorm.DB, inv models.Invitation) {
	var title, content string
	var notifType models.NotificationType

	switch inv.ResourceType {
	case models.ResourceTypeTeam:
		title = "Holdinvitation afvist"
		content = fmt.Sprintf("%s har afvist din holdinvitation", inv.Invitee.FirstName)
		notifType = models.NotifTypeTeamDecline

	case models.ResourceTypeFriend:
		title = "Venneanmodning afvist"
		content = fmt.Sprintf("%s har afvist din venneanmodning", inv.Invitee.FirstName)
		notifType = models.NotifTypeFriendDecline

	case models.ResourceTypeChallenge:
		title = "Udfordringsinvitation afvist"
		content = fmt.Sprintf("%s har afvist din challenge invitation", inv.Invitee.FirstName)
		notifType = models.NotifTypeChallengeDecline

	default:
		slog.Warn("Notification skipped: unknown resource type", slog.String("resource_type", string(inv.ResourceType)))
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
		InvitationID: &inv.ID,
	})
}

// ------ TEAMS ----- \\

// User removed from team
func CreateRemovedUserFromTeamNotification(db *gorm.DB, userID uint, team models.Team) {
	var title, content string

	title = "Du er blevet fjernet fra et hold"
	content = fmt.Sprintf("Du er blevet fjernet fra '%s'", team.Name)

	CreateNotification(db, NotificationParams{
		RecipientID: userID,
		Type:        models.NotifTypeTeamRemovedUser,
		Title:       title,
		Content:     content,
	})
}

// User left team, notifies the creator only
func CreateUserLeftTeamNotification(db *gorm.DB, leaver models.User, team models.Team) {
	var title, content string

	title = "En bruger har forladt holdet"
	content = fmt.Sprintf("%s har forladt '%s'", leaver.FirstName, team.Name)

	CreateNotification(db, NotificationParams{
		RecipientID: team.CreatorID,
		Type:        models.NotifTypeTeamUserLeft,
		Title:       title,
		Content:     content,
		ActorID:     &leaver.ID,
	})
}

func CreateTeamDeletedNotification(db *gorm.DB, user models.User, team models.Team) {
	var title, content string

	title = "Et hold du er medlem af er blevet slettet"
	content = fmt.Sprintf("'%s' er blevet slettet", team.Name)

	CreateNotification(db, NotificationParams{
		RecipientID: user.ID,
		Type:        models.NotifTypeTeamDeleted,
		Title:       title,
		Content:     content,
	})
}

// ------ CHALLENGES ----- \\
func CreateUserJoinedChallengeNotificationToCreator(db *gorm.DB, user models.User, challenge models.Challenge) {
	var title, content string
	title = "Ny deltager har tilmeldt sig din udfordring"
	content = fmt.Sprintf("%s har tilmeldt sig '%s'", user.FirstName, challenge.Name)

	CreateNotification(db, NotificationParams{
		RecipientID: challenge.CreatorID,
		Type:        models.NotifTypeChallengeReq,
		Title:       title,
		Content:     content,
	})
}

func CreateUserJoinedChallengeNotification(db *gorm.DB, user models.User, challenge models.Challenge) {
	var title, content string
	title = "Du har tilmeldt dig en udfordring"
	content = fmt.Sprintf("Du har tilmeldt dig '%s'", challenge.Name)

	CreateNotification(db, NotificationParams{
		RecipientID: user.ID,
		Type:        models.NotifTypeChallengeJoin,
		Title:       title,
		Content:     content,
	})
}

func CreateUserLeftChallengeNotification(db *gorm.DB, user models.User, challenge models.Challenge) {
	var title, content string
	title = "En deltager har forladt udfordringen"
	content = fmt.Sprintf("%s har forladt '%s'", user.FirstName, challenge.Name)

	CreateNotification(db, NotificationParams{
		RecipientID: challenge.CreatorID,
		Type:        models.NotifTypeChallengeUserLeft,
		Title:       title,
		Content:     content,
	})
}

func CreateNotificationUpcomingChallenge(db *gorm.DB, user models.User, challenge models.Challenge, notifType models.NotificationType) {
	title := "Din udfordring starter snart!"
	content := fmt.Sprintf("Din udfordring kl. '%s' starter snart", challenge.StartTime.Format("15:04"))
	rType := "challenge"
	CreateNotification(db, NotificationParams{
		RecipientID:  user.ID,
		Type:         notifType,
		Title:        title,
		Content:      content,
		ResourceID:   &challenge.ID,
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
