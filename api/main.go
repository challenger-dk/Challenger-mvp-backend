package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"server/common/config"
	"server/common/logger" // Import the logger package

	"server/api/cron"
	"server/api/routes"

	"github.com/go-chi/chi/v5"
)

func main() {
	// 1. Initialize Global Logger first
	logger.InitLogger()
	slog.Info("🚀 Initializing Challenger Backend...")

	config.LoadConfig()
	config.ConnectDatabase()

	// Run Atlas migrations
	// By default, migrations are NOT run automatically for better control
	// Set AUTO_MIGRATE=true to enable auto-migration (useful for local development)
	// In production, run migrations manually using: make migrate-up
	autoMigrate := os.Getenv("AUTO_MIGRATE")
	if autoMigrate == "true" {
		slog.Info("🔄 Auto-migration enabled (AUTO_MIGRATE=true)")
		if err := config.RunAtlasMigrations(); err != nil {
			slog.Error("Failed to run migrations", "error", err)
			os.Exit(1)
		}
	} else {
		slog.Info("⏭️  Auto-migration disabled. Run 'make migrate-up' to apply migrations")
	}

	// Seed sports data if not already seeded
	if err := config.SeedSports(); err != nil {
		slog.Error("Failed to seed sports", "error", err)
		os.Exit(1)
	}

	cron.Start()

	r := chi.NewRouter()
	routes.RegisterRoutes(r)

	// Port from Cloud Run environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	slog.Info("Starting server on :" + port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
