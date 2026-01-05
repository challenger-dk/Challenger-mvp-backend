package main

import (
	"log"
	"log/slog"
	"os"

	"server/common/config"
)

func main() {
	slog.Info("🔄 Starting migration job...")

	// Load configuration
	config.LoadConfig()

	// Connect to database
	config.ConnectDatabase()

	// Run migrations
	if err := config.RunAtlasMigrations(); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	slog.Info("✅ Migrations completed successfully")
	os.Exit(0)
}
