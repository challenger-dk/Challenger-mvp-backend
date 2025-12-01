package services

import (
	"errors"
	"fmt"
	"server/common/appError"
	"server/common/config"
	"server/common/models"
	commonServices "server/common/services"

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
	if commonServices.IsBlocked(invitation.InviteeId, invitation.InviterId) {
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
			commonServices.CreateInvitationNotification(tx, *invitation)

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
			return appError.ErrInvitationAccepted

		case models.StatusDeclined:
			// Resend by setting status back to pending
			err := tx.Model(&existing).
				Update("status", models.StatusPending).
				Error
			if err != nil {
				return err
			}
			// Re-notify
			commonServices.CreateInvitationNotification(tx, existing)
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
			commonServices.CreateAcceptedInvitationNotification(tx, invitation)

		case models.ResourceTypeFriend:
			err = createFriendship(invitation.InviterId, invitation.InviteeId, tx)
			if err != nil {
				return err
			}

			// Send notification
			commonServices.CreateAcceptedInvitationNotification(tx, invitation)

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
			commonServices.CreateAcceptedInvitationNotification(tx, invitation)

		default:
			return appError.ErrUnknownResource
		}

		invitation.Status = models.StatusAccepted
		err = tx.Save(&invitation).Error
		if err != nil {
			return err
		}

		// Mark the original invitation notification as irrelevant
		commonServices.HideNotificationByInvitationID(invitationId)

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
		commonServices.CreateDeclinedInvitationNotification(tx, invitation)

		// Mark the original invitation notification as irrelevant
		commonServices.HideNotificationByInvitationID(invitationId)

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

// addUserToChallenge adds a user to a challenge
func addUserToChallenge(challengeId uint, userId uint, db *gorm.DB) error {
	var c models.Challenge
	var u models.User

	err := db.First(&c, challengeId).Error
	if err != nil {
		return err
	}

	err = db.First(&u, userId).Error
	if err != nil {
		return err
	}

	err = db.Model(&c).
		Association("Users").
		Append(&u)

	if err != nil {
		return err
	}

	return nil
}
