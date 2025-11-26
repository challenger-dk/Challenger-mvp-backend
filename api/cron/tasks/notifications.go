package tasks

import (
	"log"
	"server/common/config"
	"server/common/models"
	"time"

	"gorm.io/gorm"
)

// RunCleanupNotifications is the entry point for the cron job.
// It handles configuration (e.g. daysOld) and logging, then calls the logic function.
// This keeps the cron schedule registry clean.
func RunCleanupNotifications() {
	log.Println("⏰ Cron: Starting scheduled cleanup...")

	// Configuration: Delete notifications older than 30 days
	const daysOld = 30

	err := deleteOldIrrelevantNotifications(daysOld)
	if err != nil {
		log.Printf("❌ Cron: Error running cleanup: %v", err)
	} else {
		log.Println("✅ Cron: Cleanup completed successfully")
	}
}

// DeleteIrrelevantNotifications permanently removes notifications that are marked irrelevant
// and are older than the specified number of days.
func deleteOldIrrelevantNotifications(daysOld int) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Calculate the cutoff date
		cutoff := time.Now().AddDate(0, 0, -daysOld)

		// Perform the delete query using the transaction handle 'tx'
		result := tx.Where("is_relevant = ? AND created_at < ?", false, cutoff).
			Delete(&models.Notification{})

		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}
