package cron

import (
	"log"
	"os"
	"server/api/cron/tasks"
	"server/common/config"

	"github.com/robfig/cron/v3"
)

func Start() {
	// 1. Check if Cron is enabled in this environment
	// This prevents multiple containers from running the same job simultaneously.
	if !config.AppConfig.EnableCron {
		log.Println("🛑 Cron scheduler disabled on this instance")
		return
	}

	// 2. Initialize with Safety Middleware
	// - Recover: Captures panics so the server doesn't crash if a job fails.
	// - SkipIfStillRunning: Prevents the same job from piling up if the previous one is slow.
	logger := cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))

	c := cron.New(
		cron.WithChain(
			cron.Recover(logger),
			cron.SkipIfStillRunning(logger),
		),
	)

	// 3. Register Tasks
	addTasks(c)

	// 4. Start
	c.Start()
	log.Println("✅ Cron scheduler started (Recover & SkipOverlap enabled)")
}

// AddTasks manages the registration of all background jobs.
// To add a new job, simply add a new line here.
func addTasks(c *cron.Cron) {
	var err error

	// --- Job 1: Notification Cleanup ---
	// Schedule: Run every day at midnight ("0 0 * * *")
	_, err = c.AddFunc("0 0 * * *", tasks.RunCleanupNotifications)
	if err != nil {
		log.Fatalf("Error scheduling CleanupNotifications: %v", err)
	}
}
