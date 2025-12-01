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
		// Pretty text for local development
		handler = NewConsoleHandler(os.Stdout)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
