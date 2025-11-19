package services

import (
	"errors"
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
			return tx.Model(&existing).
				Update("status", models.StatusPending).
				Error
		}

		return appError.ErrUnhandledInvitationStatus
	})
}

func AcceptInvitation(invitationId uint, currentUserId uint) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var invitation models.Invitation

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
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

		case models.ResourceTypeFriend:
			err = createFriendship(invitation.InviterId, invitation.InviteeId, tx)
			if err != nil {
				return err
			}

		default:
			return appError.ErrUnknownResource
		}

		invitation.Status = models.StatusAccepted
		err = tx.Save(&invitation).Error
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func DeclineInvitation(invitationId uint, currentUserId uint) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var invitation models.Invitation

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
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

		return nil
	})

	return err
}

// Private

// getResource fetches the resource associated with an invitation
func getResource(invitation models.Invitation, db *gorm.DB) (any, error) {
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

	default:
		return nil, appError.ErrUnknownResource
	}
}
