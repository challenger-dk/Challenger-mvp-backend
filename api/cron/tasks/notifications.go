package tasks

import (
	"log/slog"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

// ------- RUNNERS ------- \\

func RunCleanupNotifications() {
	slog.Info("⏰ Cron: Starting notifications cleanup...")

	const daysOld = 14

	err := deleteOldIrrelevantNotifications(daysOld)
	if err != nil {
		slog.Error("❌ Cron: Error running notifications cleanup", "error", err)
	} else {
		slog.Info("✅ Cron: Notifications cleanup completed successfully")
	}
}

// ------- TASKS ------- \\

func deleteOldIrrelevantNotifications(daysOld int) error {
	cutoff := NowFunc().AddDate(0, 0, -daysOld)

	return config.DB.Transaction(func(tx *gorm.DB) error {
		return tx.
			Where("is_relevant = ? AND created_at < ?", false, cutoff).
			Delete(&models.Notification{}).Error
	})
}
