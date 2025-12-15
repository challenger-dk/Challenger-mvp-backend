package tasks

import (
	"log/slog"
	"server/common/config"
	"server/common/models"
	"time"

	"gorm.io/gorm"
)

func RunCleanupNotifications() {
	slog.Info("⏰ Cron: Starting scheduled cleanup...")

	const daysOld = 30

	err := deleteOldIrrelevantNotifications(daysOld)
	if err != nil {
		slog.Error("❌ Cron: Error running cleanup", "error", err)
	} else {
		slog.Info("✅ Cron: Cleanup completed successfully")
	}
}

func deleteOldIrrelevantNotifications(daysOld int) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		cutoff := time.Now().AddDate(0, 0, -daysOld)
		return tx.Where("is_relevant = ? AND created_at < ?", false, cutoff).
			Delete(&models.Notification{}).Error
	})
}

func RunUpdateExpiredChallenges() {
	slog.Info("⏰ Cron: Starting expired challenges update...")

	err := updateExpiredChallenges()
	if err != nil {
		slog.Error("❌ Cron: Error updating expired challenges", "error", err)
	} else {
		slog.Info("✅ Cron: Expired challenges update completed successfully")
	}
}

func updateExpiredChallenges() error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
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
