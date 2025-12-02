package services

import (
	"errors"
	"fmt"
	"server/common/appError"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// --- GET ---
func GetInvitationsByUserId(id uint) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := config.DB.Preload("Inviter").
		Where("invitee_id = ?", id).
		Find(&invitations).
		Error

	if err != nil {
		return nil, err
	}

	return invitations, nil
}

// --- POST ---
func SendInvitation(invitation *models.Invitation) error {
	if invitation.InviterId == invitation.InviteeId {
		return appError.ErrInviteSameUser
	}

	// Check if Invitee has blocked Inviter
	if IsBlocked(invitation.InviteeId, invitation.InviterId) {
		return appError.ErrUserBlocked
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		var existing models.Invitation

		err := tx.Where(models.Invitation{
			InviterId:    invitation.InviterId,
			InviteeId:    invitation.InviteeId,
			ResourceType: invitation.ResourceType,
			ResourceID:   invitation.ResourceID,
		}).First(&existing).Error

		// Create new invitation if none exists
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			createErr := tx.Create(invitation).Error
			if createErr != nil {
				return createErr
			}

			// Successfully send invitation.
			// Create notification
			CreateInvitationNotification(tx, *invitation)

			return nil
		}

		if err != nil {
			return err
		}

		// Invitation already exists, handle based on status
		switch existing.Status {
		case models.StatusPending:
			return appError.ErrInvitationPending

		case models.StatusAccepted:
			// Check if they are currently part of the resource (Team, Friend, Challenge)
			// If they ARE part of it, return error.
			// If they are NOT part of it (they left), reset invitation to pending.

			isActive := false
			switch existing.ResourceType {
			case models.ResourceTypeTeam:
				var count int64
				err := tx.Table("user_teams").
					Where("user_id = ? AND team_id = ?", existing.InviteeId, existing.ResourceID).
					Count(&count).Error
				if err != nil {
					return err
				}
				isActive = count > 0

			case models.ResourceTypeFriend:
				var count int64
				// Check friendship (friends are usually stored bidirectionally or checked via association)
				err := tx.Table("user_friends").
					Where("user_id = ? AND friend_id = ?", existing.InviteeId, existing.InviterId).
					Count(&count).Error
				if err != nil {
					return err
				}
				isActive = count > 0

			case models.ResourceTypeChallenge:
				var count int64
				err := tx.Table("user_challenges").
					Where("user_id = ? AND challenge_id = ?", existing.InviteeId, existing.ResourceID).
					Count(&count).Error
				if err != nil {
					return err
				}
				isActive = count > 0
			}

			// If they are still active members/friends, we cannot invite them again
			if isActive {
				return appError.ErrInvitationAccepted
			}

			// If not active, they left/unfriended. Reset invitation to pending.
			err := tx.Model(&existing).
				Update("status", models.StatusPending).
				Error
			if err != nil {
				return err
			}

			// Re-notify
			CreateInvitationNotification(tx, existing)
			return nil

		case models.StatusDeclined:
			// Resend by setting status back to pending
			err := tx.Model(&existing).
				Update("status", models.StatusPending).
				Error
			if err != nil {
				return err
			}
			// Re-notify
			CreateInvitationNotification(tx, existing)
			return nil
		}

		return appError.ErrUnhandledInvitationStatus
	})
}

func AcceptInvitation(invitationId uint, currentUserId uint) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var invitation models.Invitation

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Inviter").
			Preload("Invitee").
			First(&invitation, invitationId).
			Error

		if err != nil {
			return err
		}

		if invitation.InviteeId != currentUserId {
			return appError.ErrUnauthorized
		}

		if invitation.Status != models.StatusPending {
			return appError.ErrInvitationProcessed
		}

		switch invitation.ResourceType {
		case models.ResourceTypeTeam:
			// Find resource
			resource, err := getResource(invitation, tx)
			if err != nil {
				return err
			}

			team, ok := resource.(models.Team)
			if !ok {
				return appError.ErrServerError
			}

			err = addUserToTeam(team.ID, invitation.InviteeId, tx)
			if err != nil {
				return err
			}

			// Send notification
			CreateAcceptedInvitationNotification(tx, invitation)

		case models.ResourceTypeFriend:
			err = createFriendship(invitation.InviterId, invitation.InviteeId, tx)
			if err != nil {
				return err
			}

			// Send notification
			CreateAcceptedInvitationNotification(tx, invitation)

		case models.ResourceTypeChallenge:
			// Find resource
			resource, err := getResource(invitation, tx)
			if err != nil {
				return err
			}

			challenge, ok := resource.(models.Challenge)
			if !ok {
				return appError.ErrServerError
			}

			err = addUserToChallenge(challenge.ID, invitation.InviteeId, tx)
			if err != nil {
				return err
			}

			// Send notification
			CreateAcceptedInvitationNotification(tx, invitation)

		default:
			return appError.ErrUnknownResource
		}

		invitation.Status = models.StatusAccepted
		err = tx.Save(&invitation).Error
		if err != nil {
			return err
		}

		// Mark the original invitation notification as irrelevant
		HideNotificationByInvitationID(invitationId)

		return nil
	})

	return err
}

func DeclineInvitation(invitationId uint, currentUserId uint) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var invitation models.Invitation

		// FIXED: Preloads must happen BEFORE First()
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Inviter").
			Preload("Invitee").
			First(&invitation, invitationId).
			Error

		if err != nil {
			return err
		}

		if invitation.InviteeId != currentUserId {
			return appError.ErrUnauthorized
		}

		if invitation.Status != models.StatusPending {
			return appError.ErrInvitationProcessed
		}

		invitation.Status = models.StatusDeclined
		err = tx.Save(&invitation).Error
		if err != nil {
			return err
		}

		// Send notification
		CreateDeclinedInvitationNotification(tx, invitation)

		// Mark the original invitation notification as irrelevant
		HideNotificationByInvitationID(invitationId)

		return nil
	})

	return err
}

// Private

// getResource fetches the resource associated with an invitation
func getResource(invitation models.Invitation, db *gorm.DB) (any, error) {
	fmt.Println("ResourceType:", invitation.ResourceType)
	fmt.Println("ResourceID:", invitation.ResourceID)
	switch invitation.ResourceType {
	case models.ResourceTypeTeam:
		var team models.Team
		err := db.Preload("Users").
			First(&team, invitation.ResourceID).
			Error

		if err != nil {
			return nil, err
		}
		return team, nil

	case models.ResourceTypeChallenge:
		var challenge models.Challenge
		err := db.Preload("Users").
			First(&challenge, invitation.ResourceID).
			Error

		if err != nil {
			return nil, err
		}
		return challenge, nil

	default:
		return nil, appError.ErrUnknownResource
	}
}
