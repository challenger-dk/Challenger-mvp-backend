package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"server/common/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		AppConfig.DBHost, AppConfig.DBUser, AppConfig.DBPassword, AppConfig.DBName, AppConfig.DBPort)

	maxRetries := 10
	retryDelay := 2 * time.Second

	var database *gorm.DB
	var err error

	for i := range maxRetries {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			DB = database
			slog.Info("✅ Connected to database")
			return
		}

		if i < maxRetries-1 {
			slog.Warn("Failed to initialize database, retrying...",
				"error", err,
				"attempt", i+1,
				"max_retries", maxRetries,
				"retry_delay", retryDelay,
			)
			time.Sleep(retryDelay)
			retryDelay *= 2
			if retryDelay > 10*time.Second {
				retryDelay = 10 * time.Second
			}
		}
	}

	slog.Error("Failed to connect to database after retries", "error", err)
	os.Exit(1)
}

func MigrateDB() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.Challenge{},
		&models.Sport{},
		&models.Invitation{},
		&models.Location{},
		&models.Notification{},
		&models.UserSettings{},
		&models.Message{},
		&models.Report{},
		&models.UserChat{},
		&models.Chat{},
	)

	if err != nil {
		slog.Error("Failed to migrate database", "error", err)
		os.Exit(1)
	}
	slog.Info("Database migrated successfully")
}
