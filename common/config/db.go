package config

import (
	"fmt"
	"log"
	"time"

	"server/common/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Build the DSN from the loaded AppConfig struct
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		AppConfig.DBHost, AppConfig.DBUser, AppConfig.DBPassword, AppConfig.DBName, AppConfig.DBPort)

	// Retry connection with exponential backoff
	maxRetries := 10
	retryDelay := 2 * time.Second

	var database *gorm.DB
	var err error

	for i := range maxRetries {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			DB = database
			fmt.Println("✅ Connected to database")
			return
		}

		if i < maxRetries-1 {
			log.Printf("[error] failed to initialize database, got error %v. Retrying in %v... (attempt %d/%d)", err, retryDelay, i+1, maxRetries)
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
			if retryDelay > 10*time.Second {
				retryDelay = 10 * time.Second // Cap at 10 seconds
			}
		}
	}

	log.Fatal("Failed to connect to database:", err)
}

func MigrateDB() {
	err := DB.AutoMigrate(&models.User{},
		&models.Team{},
		&models.Challenge{},
		&models.Sport{},
		&models.Invitation{},
		&models.Location{},
		&models.Notification{},
		&models.UserSettings{})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}
