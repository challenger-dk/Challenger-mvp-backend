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
