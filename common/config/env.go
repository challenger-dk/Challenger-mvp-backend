package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Database connection settings
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`

	// JWT Secret
	JWTSecret string `env:"JWT_SECRET,required"`

	// Cron Settings
	// In production, set this to "true" on only ONE instance/container
	EnableCron bool `env:"ENABLE_CRON" envDefault:"true"`
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	if err := env.Parse(&AppConfig); err != nil {
		log.Fatalf("Failed to parse config: %+v", err)
	}

	log.Println("✅ Configuration loaded")
}
