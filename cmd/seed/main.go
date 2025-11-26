package main

import (
	"log"

	"server/common/config"
	"server/common/seed"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	config.ConnectDatabase()

	// Ensure PostGIS extension is created
	err := config.DB.Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error
	if err != nil {
		log.Fatal("Failed to create PostGIS extension:", err)
	}

	// Run migrations
	config.MigrateDB()

	// Seed the database
	if err := seed.SeedDatabase(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("🎉 Seeding completed successfully!")
}
