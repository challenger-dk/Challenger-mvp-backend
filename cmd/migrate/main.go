package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"

	"server/common/config"
)

func main() {
	slog.Info("🔄 Starting migration job...")

	// Load configuration
	config.LoadConfig()

	// Connect to database
	config.ConnectDatabase()

	// Check if we should baseline (set via env var)
	if os.Getenv("BASELINE_MIGRATIONS") == "true" {
		slog.Info("🔧 Baselining migrations (marking existing migrations as applied)...")
		if err := baselineMigrations(); err != nil {
			log.Fatalf("❌ Failed to baseline migrations: %v", err)
		}
		slog.Info("✅ Migrations baselined successfully")
		os.Exit(0)
	}

	// Run migrations
	if err := config.RunAtlasMigrations(); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}

	slog.Info("✅ Migrations completed successfully")
	os.Exit(0)
}

// baselineMigrations marks all existing migrations as applied without running them
func baselineMigrations() error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	cmd := exec.Command("atlas", "migrate", "set",
		"--dir", "file://Database/migrations",
		"--url", dsn,
		"--baseline", "20251219224434_add_reports_table.sql", // Last migration file
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
