package services

import (
	"errors"
	"log/slog"
	"server/common/appError"
	"server/common/config"
	"server/common/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetChallengeByID(id uint, currentUserID uint) (models.Challenge, error) {
	var c models.Challenge

	err := config.DB.
		Scopes(ExcludeBlockedUsersOn(currentUserID, "creator_id")).
		Preload("Users", ExcludeBlockedUsers(currentUserID)).
		Preload("Teams").
		Preload("Creator").
		Preload("Location").
		First(&c, id).
		Error

	if err != nil {
		return models.Challenge{}, err
	}

	// Update status to completed if EndTime has passed
	updateChallengeStatusIfExpired(&c)

	return c, nil
}

// TODO: Some kind of pagination so we dont fetch all challenges
func GetChallenges(currentUserID uint) ([]models.Challenge, error) {
	var challenges []models.Challenge

	err := config.DB.
		Scopes(ExcludeBlockedUsersOn(currentUserID, "creator_id")).
		Preload("Users", ExcludeBlockedUsers(currentUserID)).
		Preload("Teams").
		Preload("Creator").
		Preload("Location").
		Find(&challenges).
		Error

	if err != nil {
		return nil, err
	}

	// Update status to completed for challenges where EndTime has passed
	for i := range challenges {
		updateChallengeStatusIfExpired(&challenges[i])
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

		// Find or create the location first
		location, err := FindOrCreateLocation(tx, c.Location)
		if err != nil {
			return err
		}

		// Set the LocationID and clear the Location object to avoid GORM trying to create it
		c.LocationID = location.ID
		c.Location = models.Location{}

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
				CreateInvitationNotification(tx, invitation)
			} else if err == nil {
				// Invitation already exists, handle based on status
				switch existing.Status {
				case models.StatusPending:
					// Already pending, skip
					continue
				case models.StatusAccepted:
					// Already accepted, add user to challenge (with capacity check)
					err = addUserToChallenge(c.ID, userId, tx)
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
					CreateInvitationNotification(tx, existing)
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

	// Successfully created challenge, notify creator
	rType := models.ResourceTypeChallenge
	CreateNotification(config.DB, NotificationParams{
		RecipientID:  c.CreatorID,
		Type:         models.NotifTypeChallengeCreated,
		Title:        "Udfordring oprettet",
		Content:      "Din udfordring er live og klar til deltagere.",
		ResourceID:   &c.ID,
		ResourceType: &rType,
	})

	// Create challenge conversation with initial members
	memberIDs := make([]uint, len(c.Users))
	for i, u := range c.Users {
		memberIDs[i] = u.ID
	}
	if err := SyncChallengeConversationMembers(c.ID, memberIDs); err != nil {
		// Log error but don't fail the request
		// Challenge conversation can be created later
		slog.Warn("Failed to create challenge conversation for challenge",
			slog.Int("challenge_id", int(c.ID)),
			slog.Any("error", err),
		)
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

		// Update basic fields
		if ch.Name != "" {
			c.Name = ch.Name
		}

		if ch.Description != "" {
			c.Description = ch.Description
		}

		if ch.Sport != "" {
			c.Sport = ch.Sport
		}

		// Update location - find or create if location data is provided
		if ch.Location.Address != "" {
			location, err := FindOrCreateLocation(tx, ch.Location)
			if err != nil {
				return err
			}
			c.LocationID = location.ID
		}

		// Update boolean fields (always update since they can be true/false)
		c.IsIndoor = ch.IsIndoor
		c.IsPublic = ch.IsPublic
		c.HasCost = ch.HasCost

		// Update pointer fields (string)
		if ch.Comment != nil {
			c.Comment = ch.Comment
		}

		if ch.PlayFor != nil {
			c.PlayFor = ch.PlayFor
		}

		// Update pointer fields (numeric)
		if ch.TeamSize != nil {
			c.TeamSize = ch.TeamSize
		}

		if ch.Distance != nil {
			c.Distance = ch.Distance
		}

		if ch.Participants != nil {
			c.Participants = ch.Participants
		}

		// Update time fields
		if !ch.Date.IsZero() {
			c.Date = ch.Date
		}

		if !ch.StartTime.IsZero() {
			c.StartTime = ch.StartTime
		}

		if ch.EndTime != nil {
			c.EndTime = ch.EndTime
		}

		// Update status and type
		if ch.Status != "" {
			c.Status = ch.Status
		}

		if ch.Type != "" {
			c.Type = ch.Type
		}

		return tx.Save(&c).Error
	})
}

func JoinChallenge(id uint, userId uint) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		return addUserToChallenge(id, userId, tx)
	})
	if err != nil {
		return err
	}

	// Sync challenge conversation members after successful join
	var challenge models.Challenge
	if err := config.DB.Preload("Users").First(&challenge, id).Error; err != nil {
		return err
	}

	memberIDs := make([]uint, len(challenge.Users))
	for i, u := range challenge.Users {
		memberIDs[i] = u.ID
	}

	if err := SyncChallengeConversationMembers(id, memberIDs); err != nil {
		// Log error but don't fail the request
		slog.Warn("Failed to sync challenge conversation after user joined",
			slog.Uint64("challenge_id", uint64(id)),
			slog.Uint64("user_id", uint64(userId)),
			slog.Any("error", err),
		)
	}

	return nil
}

func LeaveChallenge(id uint, userId uint) error {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var c models.Challenge
		var u models.User

		err := tx.First(&c, id).Error
		if err != nil {
			return err
		}

		err = tx.First(&u, userId).Error
		if err != nil {
			return err
		}

		return tx.Model(&c).
			Association("Users").
			Delete(&u)
	})
	if err != nil {
		return err
	}

	// Sync challenge conversation members after successful leave
	var challenge models.Challenge
	if err := config.DB.Preload("Users").First(&challenge, id).Error; err != nil {
		return err
	}

	memberIDs := make([]uint, len(challenge.Users))
	for i, u := range challenge.Users {
		memberIDs[i] = u.ID
	}

	if err := SyncChallengeConversationMembers(id, memberIDs); err != nil {
		// Log error but don't fail the request
		slog.Warn("Failed to sync challenge conversation after user left",
			slog.Uint64("challenge_id", uint64(id)),
			slog.Uint64("user_id", uint64(userId)),
			slog.Any("error", err),
		)
	}

	return nil
}

func DeleteChallenge(id uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var c models.Challenge

		if err := tx.First(&c, id).Error; err != nil {
			return err
		}

		// Soft delete: this sets c.DeletedAt, keeps row & associations
		if err := tx.Delete(&c).Error; err != nil {
			return err
		}

		return nil
	})
}

// updateChallengeStatusIfExpired checks if a challenge's EndTime has passed
// and updates the status to "completed" if it has and the challenge is not already completed
func updateChallengeStatusIfExpired(c *models.Challenge) {
	if c.EndTime == nil {
		return
	}

	now := time.Now()
	if c.EndTime.Before(now) && c.Status != models.ChallengeStatusCompleted {
		// Only update if not already completed to avoid unnecessary database writes
		config.DB.Model(c).Update("status", models.ChallengeStatusCompleted)
		c.Status = models.ChallengeStatusCompleted
	}
}

// addUserToChallenge adds a user to a challenge
func addUserToChallenge(challengeId uint, userId uint, db *gorm.DB) error {
	var c models.Challenge
	var u models.User

	// Lock the challenge row up-front to prevent race conditions
	if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&c, challengeId).Error; err != nil {
		return err
	}

	// Load the user
	if err := db.First(&u, userId).Error; err != nil {
		return err
	}

	// Check if user is already in the challenge
	var alreadyMember int64
	if err := db.Table("user_challenges").
		Where("user_id = ? AND challenge_id = ?", userId, challengeId).
		Count(&alreadyMember).Error; err != nil {
		return err
	}
	if alreadyMember > 0 {
		return appError.ErrUserAlreadyInChallenge
	}

	// Check capacity
	if c.Participants != nil {
		var currentCount int64
		if err := db.Table("user_challenges").
			Where("challenge_id = ?", challengeId).
			Count(&currentCount).Error; err != nil {
			return err
		}
		if currentCount >= int64(*c.Participants) {
			return appError.ErrChallengeFullParticipation
		}
	}

	// Append membership
	if err := db.Model(&c).
		Association("Users").
		Append(&u); err != nil {
		return err
	}

	// Re-count after insert
	var newCount int64
	if err := db.Table("user_challenges").
		Where("challenge_id = ?", challengeId).
		Count(&newCount).Error; err != nil {
		return err
	}

	isFull := c.Participants != nil && newCount == int64(*c.Participants)

	// If full, update challenge status to CONFIRMED
	if isFull && c.Status != models.ChallengeConfirmed && c.Status != models.ChallengeStatusCompleted {
		if err := db.Model(&c).
			Update("status", models.ChallengeConfirmed).Error; err != nil {
			return err
		}
		c.Status = models.ChallengeConfirmed
	}

	// Load creator for notification
	var creator models.User
	if err := db.First(&creator, c.CreatorID).Error; err != nil {
		return err
	}

	// Notify the joining user
	CreateUserJoinedChallengeNotification(db, u, c)

	// Notify creator: either challenge became full, or someone joined
	if isFull {
		CreateNotificationChallengeFullParticipation(db, creator, c)
	} else {
		CreateUserJoinedChallengeNotificationToCreator(db, u, c)
	}

	return nil
}

func ConfirmChallenge(id uint, user *models.User) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var c models.Challenge
		var u models.User

		err := tx.First(&c, id).Error
		if err != nil {
			return err
		}

		err = tx.First(&u, user.ID).Error
		if err != nil {
			return err
		}
		if c.CreatorID != user.ID {
			return appError.ErrUnauthorized
		}
		if c.Status != models.ChallengeStatusPending {
			return appError.ErrChallengeAlreadyConfirmed
		}

		return tx.Model(&c).
			Update("status", models.ChallengeConfirmed).Error
	})
}
