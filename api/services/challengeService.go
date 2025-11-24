package services

import (
	"errors"
	"server/common/config"
	"server/common/models"
	commonServices "server/common/services"

	"gorm.io/gorm"
)

func GetChallengeByID(id uint) (models.Challenge, error) {
	var c models.Challenge

	err := config.DB.Preload("Users").
		Preload("Teams").
		Preload("Creator").
		Preload("Location").
		First(&c, id).
		Error

	if err != nil {
		return models.Challenge{}, err
	}

	return c, nil
}

// TODO: Some kind of pagination so we dont fetch all challenges
func GetChallenges() ([]models.Challenge, error) {
	var challenges []models.Challenge

	err := config.DB.Preload("Users").
		Preload("Teams").
		Preload("Creator").
		Preload("Location").
		Find(&challenges).
		Error

	if err != nil {
		return nil, err
	}

	return challenges, nil
}

func CreateChallenge(c models.Challenge, invitedUserIds []uint) (models.Challenge, error) {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		creator := models.User{}

		// ensure creator exists
		err := tx.First(&creator, c.CreatorID).Error
		if err != nil {
			return err
		}

		c.CreatorID = creator.ID
		c.Creator = models.User{}

		err = tx.Create(&c).Error
		if err != nil {
			return err
		}

		// Automatically add creator to challenge Users
		err = tx.Model(&c).
			Association("Users").
			Append(&creator)
		if err != nil {
			return err
		}

		// Create invitations for each invited user
		for _, userId := range invitedUserIds {
			// Skip if trying to invite the creator (they're already added)
			if userId == creator.ID {
				continue
			}

			invitation := models.Invitation{
				InviterId:    creator.ID,
				InviteeId:    userId,
				ResourceType: models.ResourceTypeChallenge,
				ResourceID:   c.ID,
				Status:       models.StatusPending,
			}

			// Use SendInvitation logic but within the same transaction
			var existing models.Invitation
			err := tx.Where(models.Invitation{
				InviterId:    invitation.InviterId,
				InviteeId:    invitation.InviteeId,
				ResourceType: invitation.ResourceType,
				ResourceID:   invitation.ResourceID,
			}).First(&existing).Error

			// Create new invitation if none exists
			if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
				createErr := tx.Create(&invitation).Error
				if createErr != nil {
					return createErr
				}

				// Create notification
				commonServices.CreateInvitationNotification(tx, invitation)
			} else if err == nil {
				// Invitation already exists, handle based on status
				switch existing.Status {
				case models.StatusPending:
					// Already pending, skip
					continue
				case models.StatusAccepted:
					// Already accepted, add user to challenge
					var user models.User
					err = tx.First(&user, userId).Error
					if err != nil {
						return err
					}
					err = tx.Model(&c).
						Association("Users").
						Append(&user)
					if err != nil {
						return err
					}
				case models.StatusDeclined:
					// Resend by setting status back to pending
					err = tx.Model(&existing).
						Update("status", models.StatusPending).
						Error
					if err != nil {
						return err
					}
					commonServices.CreateInvitationNotification(tx, existing)
				}
			} else {
				// Some other error occurred
				return err
			}
		}

		err = tx.Preload("Users").
			Preload("Teams").
			Preload("Creator").
			Preload("Location").
			First(&c, c.ID).
			Error

		return err
	})

	if err != nil {
		return models.Challenge{}, err
	}

	return c, nil
}

func UpdateChallenge(id uint, ch models.Challenge) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var c models.Challenge

		err := tx.First(&c, id).Error
		if err != nil {
			return err
		}

		if ch.Name != "" {
			c.Name = ch.Name
		}

		if ch.Description != "" {
			c.Description = ch.Description
		}

		if ch.Sport != "" {
			c.Sport = ch.Sport
		}

		if ch.Location.ID != 0 {
			c.LocationID = ch.Location.ID
		}

		if ch.TeamSize != nil {
			c.TeamSize = ch.TeamSize
		}

		return tx.Save(&c).Error
	})
}

func DeleteChallenge(id uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var c models.Challenge

		err := tx.First(&c, id).Error
		if err != nil {
			return err
		}

		// clear many2many associations to avoid orphan join rows
		err = tx.Model(&c).
			Association("Teams").
			Clear()

		if err != nil {
			return err
		}

		err = tx.Model(&c).
			Association("Users").
			Clear()

		if err != nil {
			return err
		}

		return tx.Delete(&c).Error
	})
}
