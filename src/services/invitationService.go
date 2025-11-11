package services

import (
	"errors"
	"server/config"
	"server/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// --- GET ---
func GetInvitationsByUserId(id uint) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := config.DB.Preload("Inviter").Where("invitee_id = ?", id).Find(&invitations).Error
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

// --- POST ---
func SendInvitation(invitation *models.Invitation) error {
	if invitation.InviterId == invitation.InviteeId {
		return errors.New("inviter and invitee cannot be the same user")
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
			if createErr := tx.Create(invitation).Error; createErr != nil {
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
			return errors.New("invitation is already pending")
		case models.StatusAccepted:
			return errors.New("user has already accepted this invitation")
		case models.StatusDeclined:
			// Resend by setting status back to pending
			return tx.Model(&existing).Update("status", models.StatusPending).Error
		}

		return errors.New("unhandled invitation status")
	})
}

func AcceptInvitation(invitationId uint) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var invitation models.Invitation

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&invitation, invitationId).Error
		if err != nil {
			return err
		}

		if invitation.Status != models.StatusPending {
			return errors.New("invitation already processed")
		}

		// Find resource
		resource, err := getResource(invitation, tx)
		if err != nil {
			return err
		}

		// Add user to resource
		switch res := resource.(type) {
		case models.Team:
			err = addUserToTeam(res.ID, invitation.InviteeId, tx)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown resource type")
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

func DeclineInvitation(invitationId uint) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var invitation models.Invitation

		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&invitation, invitationId).Error
		if err != nil {
			return err
		}

		if invitation.Status != models.StatusPending {
			return errors.New("invitation already processed")
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
func getResource(invitation models.Invitation, db *gorm.DB) (any, error) {
	switch invitation.ResourceType {
	case models.ResourceTypeTeam:
		var team models.Team
		err := db.Preload("Users").First(&team, invitation.ResourceID).Error
		if err != nil {
			return nil, err
		}
		return team, nil
	default:
		return nil, errors.New("unknown resource type")
	}
}
