package cron

import (
	"log/slog"
	"os"
	"server/api/cron/tasks"
	"server/common/config"

	"github.com/robfig/cron/v3"
)

// cronLogger adapts slog to the cron.Logger interface
type cronLogger struct{}

func (l cronLogger) Info(msg string, keysAndValues ...any) {
	slog.Info(msg, keysAndValues...)
}

func (l cronLogger) Error(err error, msg string, keysAndValues ...any) {
	args := append([]any{"error", err}, keysAndValues...)
	slog.Error(msg, args...)
}

func Start() {
	if !config.AppConfig.EnableCron {
		slog.Info("🛑 Cron scheduler disabled on this instance")
		return
	}

	// Use custom logger adapter
	c := cron.New(
		cron.WithChain(
			cron.Recover(cronLogger{}),
			cron.SkipIfStillRunning(cronLogger{}),
		),
		cron.WithLogger(cronLogger{}),
	)

	addTasks(c)

	c.Start()
	slog.Info("✅ Cron scheduler started", "middleware", "Recover & SkipOverlap")
}

func addTasks(c *cron.Cron) {
	_, err := c.AddFunc("0 0 * * *", tasks.RunCleanupNotifications)
	if err != nil {
		slog.Error("Error scheduling CleanupNotifications", "error", err)
		os.Exit(1)
	}

	// Run every hour to update expired challenges
	_, err = c.AddFunc("0 * * * *", tasks.RunUpdateExpiredChallenges)
	if err != nil {
		slog.Error("Error scheduling UpdateExpiredChallenges", "error", err)
		os.Exit(1)
	}
}
