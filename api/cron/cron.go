package cron

import (
	"log/slog"
	"os"
	"server/api/cron/tasks"
	"server/common/config"
	"time"

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

	// Make sure to set the correct timezone for scheduling
	loc, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		// Fallback to UTC if tz data is missing in the container
		slog.Warn("Failed to load timezone, falling back to UTC", "error", err)
		loc = time.UTC
	}

	c := cron.New(
		cron.WithLocation(loc),
		cron.WithChain(
			cron.Recover(cronLogger{}),
			cron.SkipIfStillRunning(cronLogger{}),
		),
		cron.WithLogger(cronLogger{}),
	)

	addTasks(c)

	c.Start()
	slog.Info("✅ Cron scheduler started",
		"middleware", "Recover & SkipOverlap",
		"location", loc.String(),
	)
}

func addTasks(c *cron.Cron) {

	// ------- CLEANUP TASKS ------- \\

	// Run every day at midnight to cleanup old irrelevant notifications
	_, err := c.AddFunc("@daily", tasks.RunCleanupNotifications)
	if err != nil {
		slog.Error("Error scheduling RunCleanupNotifications", "error", err)
		os.Exit(1)
	}

	// Run every hour to update expired challenges
	_, err = c.AddFunc("@hourly", tasks.RunUpdateExpiredChallenges)
	if err != nil {
		slog.Error("Error scheduling RunUpdateExpiredChallenges", "error", err)
		os.Exit(1)
	}

	// ------- NOTIFI USER TASKS ------- \\

	// Notify users 24 hours before challenge start
	_, err = c.AddFunc("@every 10m", tasks.RunNotifiUserUpcommingChallenges24H)
	if err != nil {
		slog.Error("Error scheduling RunNotifiUserUpcommingChallenges24H", "error", err)
		os.Exit(1)
	}

	// Notify users 1 hour before challenge start
	_, err = c.AddFunc("@every 10m", tasks.RunNotifiUserUpcommingChallenges1H)
	if err != nil {
		slog.Error("Error scheduling RunNotifiUserUpcommingChallenges1H", "error", err)
		os.Exit(1)
	}

	// Notify users 24 hours before challenge start if they haven't answered invitation
	_, err = c.AddFunc("@every 10m", tasks.RunNotifiUserInvitedToChallengeNotAnswered24H)
	if err != nil {
		slog.Error("Error scheduling RunNotifiUserInvitedToChallengeNotAnswered24H", "error", err)
		os.Exit(1)
	}

	// Notify creators 12 hours before challenge start if participants are missing
	_, err = c.AddFunc("@every 10m", tasks.RunNotifiUserMissingParticipantsInChallenges12H)
	if err != nil {
		slog.Error("Error scheduling RunNotifiUserMissingParticipantsInChallenges12H", "error", err)
		os.Exit(1)
	}
}
