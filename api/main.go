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

	// Ensure PostGIS extension is created
	err := config.DB.Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error
	if err != nil {
		slog.Error("Failed to create PostGIS extension", "error", err)
		os.Exit(1)
	}
	slog.Info("✅ PostGIS extension ensured")

	// Run migrations (for local development with Docker)
	// In production, migrations are handled by Atlas
	if config.AppConfig.AppEnv == "development" {
		config.MigrateDB()
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
