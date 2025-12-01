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
