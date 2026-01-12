package tasks

import (
	"log/slog"
	"server/common/config"
	"server/common/models"
	"server/common/services"
	"time"

	"gorm.io/gorm"
)

func RunUpdateExpiredChallenges() {
	slog.Info("⏰ Cron: Starting expired challenges update...")

	err := updateExpiredChallenges()
	if err != nil {
		slog.Error("❌ Cron: Error updating expired challenges", "error", err)
	} else {
		slog.Info("✅ Cron: Expired challenges update completed successfully")
	}
}

// TODO:
// Notifi creators of challenges starting in ~12 hours with participants missing
/*
func RunNotifiMissingParticipantsInChallenges() {
	now := time.Now()
	windowStart := now.Add(12*time.Hour - 5*time.Minute)
	windowEnd := now.Add(12*time.Hour + 5*time.Minute)

	var challenges []models.Challenge

}
*/

// ------- RUNNERS ------- \\

func RunNotifiUserUpcommingChallenges24H() {
	slog.Info("⏰ Cron: Starting upcomming challenges notification 24h...")
	err := notifiUserUpcommingChallenges24H()
	if err != nil {
		slog.Error("❌ Cron: Error notifying users of upcomming challenges", "error", err)
	} else {
		slog.Info("✅ Cron: Upcomming challenges notification completed successfully")
	}
}

func RunNotifiUserUpcommingChallenges1H() {
	slog.Info("⏰ Cron: Starting upcomming challenges notification 1h...")
	err := notifiUserUpcommingChallenges1H()
	if err != nil {
		slog.Error("❌ Cron: Error notifying users of upcomming challenges", "error", err)
	} else {
		slog.Info("✅ Cron: Upcomming challenges notification completed successfully")
	}
}

func RunNotifiUserInvitedToChallengeNotAnswered24H() {
	slog.Info("⏰ Cron: Starting pending challenge invitation reminder (24h before start)...")
	err := notifiUserInvitedToChallengeNotAnswered24H()
	if err != nil {
		slog.Error("❌ Cron: Error notifying users about pending invitations", "error", err)
	} else {
		slog.Info("✅ Cron: Pending invitation reminder completed successfully")
	}
}

func RunNotifiUserMissingParticipantsInChallenges12H() {
	slog.Info("⏰ Cron: Starting missing participants notification (12h before start)...")
	err := notifiUserMissingParticipantsInChallenges12H()
	if err != nil {
		slog.Error("❌ Cron: Error notifying users about missing participants", "error", err)
	} else {
		slog.Info("✅ Cron: Missing participants notification completed successfully")
	}
}

// ------- IMPLEMENTATION ------- \\

func notifiUserUpcommingChallenges24H() error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		now := NowFunc()
		windowStart := now.Add(24*time.Hour - 5*time.Minute)
		windowEnd := now.Add(24*time.Hour + 5*time.Minute)

		var challenges []models.Challenge
		if err := tx.Preload("Users").
			Where("start_time BETWEEN ? AND ?", windowStart, windowEnd).
			Find(&challenges).Error; err != nil {
			return err
		}

		for _, ch := range challenges {
			for _, u := range ch.Users {
				// Skip if already exists
				var existing models.Notification
				if err := tx.Where("user_id = ? AND resource_id = ? AND type = ?",
					u.ID, ch.ID, models.NotifTypeChallengeUpcomming24H).First(&existing).Error; err == nil {
					continue
				}

				services.CreateNotificationUpcomingChallenge(tx, u, ch, models.NotifTypeChallengeUpcomming24H)
			}
		}

		return nil
	})
}

func notifiUserUpcommingChallenges1H() error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		now := NowFunc()
		windowStart := now.Add(1*time.Hour - 5*time.Minute)
		windowEnd := now.Add(1*time.Hour + 5*time.Minute)
		var challenges []models.Challenge
		if err := tx.Preload("Users").
			Where("start_time BETWEEN ? AND ?", windowStart, windowEnd).
			Find(&challenges).Error; err != nil {
			return err
		}

		for _, ch := range challenges {
			for _, u := range ch.Users {
				// Skip if already exists
				var existing models.Notification
				if err := tx.Where("user_id = ? AND resource_id = ? AND type = ?",
					u.ID, ch.ID, models.NotifTypeChallengeUpcomming1H).First(&existing).Error; err == nil {
					continue
				}

				services.CreateNotificationUpcomingChallenge(tx, u, ch, models.NotifTypeChallengeUpcomming1H)
			}
		}

		return nil
	})
}

func notifiUserInvitedToChallengeNotAnswered24H() error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		now := NowFunc()
		windowStart := now.Add(24*time.Hour - 5*time.Minute)
		windowEnd := now.Add(24*time.Hour + 5*time.Minute)

		// Fetch pending challenge invitations where the challenge starts in ~24h
		type row struct {
			models.Invitation `gorm:"embedded"` // Embed invitation model for easy access to fields (e.g. Inviter, Invitee) (for deduplication)
			StartTime         time.Time         `gorm:"column:start_time"` // Challenge start time (for deduplication)
		}

		var invitations []row
		if err := tx.
			Table("invitations").
			Joins("JOIN challenges ON challenges.id = invitations.resource_id").
			Where("invitations.resource_type = ?", models.ResourceTypeChallenge).
			Where("invitations.status = ?", models.StatusPending).
			Where("challenges.start_time BETWEEN ? AND ?", windowStart, windowEnd).
			Preload("Invitee").
			Preload("Inviter").
			Select("invitations.*, challenges.start_time").
			Find(&invitations).Error; err != nil {
			return err
		}

		if len(invitations) == 0 {
			return nil
		}

		// To avoid N+1 loading the same challenge repeatedly
		challengeByID := make(map[uint]models.Challenge, 32)

		for _, inv := range invitations {
			// Load challenge once (within transaction)
			ch, ok := challengeByID[inv.ResourceID]
			if !ok {
				var c models.Challenge
				if err := tx.First(&c, inv.ResourceID).Error; err != nil {
					// If challenge was deleted or missing, just skip this invitation
					// (alternatively: return err if you want strictness)
					continue
				}
				challengeByID[inv.ResourceID] = c
				ch = c
			}

			// Deduplicate: skip if notification already exists
			var existing models.Notification
			if err := tx.
				Where("user_id = ? AND resource_id = ? AND type = ?",
					inv.InviteeId, inv.ResourceID, models.NotifTypeChallengeNotAnswered24H).
				First(&existing).Error; err == nil {
				continue
			}

			// Ensure we have a valid invitee user to notify
			var invitee models.User
			if inv.Invitee.ID != 0 {
				invitee = inv.Invitee
			} else {
				if err := tx.First(&invitee, inv.InviteeId).Error; err != nil {
					// Missing invitee (deleted user?) — skip
					continue
				}
			}

			// Create notification to invitee
			services.CreateNotificationChallengeNotAnswered24H(tx, invitee, ch)
		}

		return nil
	})
}

func updateExpiredChallenges() error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		now := NowFunc()
		result := tx.Model(&models.Challenge{}).
			Where("end_time IS NOT NULL AND end_time < ? AND status != ?", now, models.ChallengeStatusCompleted).
			Update("status", models.ChallengeStatusCompleted)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected > 0 {
			slog.Info("✅ Cron: Updated expired challenges", "count", result.RowsAffected)
		}

		return nil
	})
}

// Notify creators of challenges starting in 12 hours with 1-2 participants missing
func notifiUserMissingParticipantsInChallenges12H() error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		now := NowFunc()
		windowStart := now.Add(12*time.Hour - 5*time.Minute)
		windowEnd := now.Add(12*time.Hour + 5*time.Minute)

		var challenges []models.Challenge
		if err := tx.Preload("Users").
			Where("start_time BETWEEN ? AND ?", windowStart, windowEnd).
			Find(&challenges).Error; err != nil {
			return err
		}

		for _, ch := range challenges {
			// Must have a participants limit to evaluate missing participants
			if ch.Participants == nil {
				continue
			}

			participantCount := len(ch.Users)
			totalNeeded := int(*ch.Participants)

			// Already full or over-capacity -> skip
			if participantCount >= totalNeeded {
				continue
			}

			// Notify creator if exists (any missing participants)
			var creator models.User
			if err := tx.First(&creator, ch.CreatorID).Error; err != nil {
				continue
			}

			// Skip if notification already exists
			var existing models.Notification
			if err := tx.Where("user_id = ? AND resource_id = ? AND type = ?",
				creator.ID, ch.ID, models.NotifTypeChallengeMissingParticipants).First(&existing).Error; err == nil {
				continue
			}

			// Send notification
			services.CreateNotificationChallengeMissingParticipants(tx, creator, ch)
		}

		return nil
	})
}
