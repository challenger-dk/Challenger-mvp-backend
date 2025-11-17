package config

import (
	"fmt"
	"log"
	"server/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Build the DSN from the loaded AppConfig struct
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		AppConfig.DBHost, AppConfig.DBUser, AppConfig.DBPassword, AppConfig.DBName, AppConfig.DBPort)

	// Open database connection
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = database
	fmt.Println("✅ Connected to database")
}

func MigrateDB() {
	err := DB.AutoMigrate(&models.User{},
		&models.Team{},
		&models.Challenge{},
		&models.Sport{},
		&models.Invitation{},
		&models.Location{})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}
