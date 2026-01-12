package logger

import (
	"log/slog"
	"os"
	"strings"
)

func InitLogger() {
	var handler slog.Handler

	// Read environment directly to avoid circular dependency with config package
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // Default to dev
	}

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug, // Show everything in dev
	}

	if strings.ToLower(env) == "production" {
		// JSON for production tools (Datadog/CloudWatch)
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Use slog's built-in text handler in dev so HandlerOptions (levels, time formatting)
		// are respected by default and reliably printed to stdout.
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Quick health-check log to confirm logger is active
	slog.Info("Logger initialized", "env", env)
}
