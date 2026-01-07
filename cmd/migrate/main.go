package main

import (
	"log"
	"log/slog"
	"os"

	"server/common/config"
)

func main() {
	slog.Info("🔄 Starting migration job...")

	// Debug: List migration files in the Docker image
	slog.Info("📁 Listing migration files in Docker image...")
	files, err := os.ReadDir("Database/migrations")
	if err != nil {
		slog.Error("Failed to read migrations directory", "error", err)
	} else {
		for _, file := range files {
			slog.Info("  Found migration file", "name", file.Name())
		}
	}

	// Load configuration
	config.LoadConfig()

	// Connect to database
	config.ConnectDatabase()

	// Check if we should reset migration history (use with caution!)
	if os.Getenv("RESET_MIGRATION_HISTORY") == "true" {
		slog.Warn("⚠️  RESET_MIGRATION_HISTORY=true - Dropping atlas_schema_revisions table...")
		if err := resetMigrationHistory(); err != nil {
			log.Fatalf("❌ Failed to reset migration history: %v", err)
		}
		slog.Info("✅ Migration history reset successfully")
	}

	// Run migrations
	if err := config.RunAtlasMigrations(); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	slog.Info("✅ Migrations completed successfully")
	os.Exit(0)
}

func resetMigrationHistory() error {
	db, err := config.DB.DB()
	if err != nil {
		return err
	}

	_, err = db.Exec("DROP TABLE IF EXISTS atlas_schema_revisions CASCADE")
	return err
}
