package main

import (
	"log"
	"net/http"
	"time"

	"server/common/config"

	"server/api/routes"

	"server/api/cron"

	"github.com/go-chi/chi/v5"
)

func main() {

	// Loads config from .env and environment variables
	config.LoadConfig()
	config.ConnectDatabase()

	// Ensure PostGIS extension is created
	err := config.DB.Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error
	if err != nil {
		log.Fatal("Failed to create PostGIS extension:", err)
	}

	config.MigrateDB()

	// Seed allowed sports
	if err := config.SeedSports(); err != nil {
		log.Fatal("Failed to seed sports:", err)
	}

	cron.Start()

	r := chi.NewRouter()
	routes.RegisterRoutes(r)

	server := &http.Server{
		Addr:         ":8000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Println("Starting server on :8000")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
