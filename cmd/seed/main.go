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

	// Run Atlas migrations
	if err := config.RunAtlasMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Seed the database
	if err := seed.SeedDatabase(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("🎉 Seeding completed successfully!")
}
