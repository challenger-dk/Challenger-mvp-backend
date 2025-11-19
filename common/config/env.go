package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
// The `env` tags are used to parse values from the environment
type Config struct {
	// Database connection settings
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`

	// JWT Secret
	JWTSecret string `env:"JWT_SECRET,required"`
}

var AppConfig Config

// LoadConfig loads configuration from .env and environment variables
// This function should be called ONCE at startup.
func LoadConfig() {
	// Load .env file (for local development)
	// This is ignored in production/CI where .env doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	// Parse the environment variables into the AppConfig struct
	if err := env.Parse(&AppConfig); err != nil {
		log.Fatalf("Failed to parse config: %+v", err)
	}

	log.Println("✅ Configuration loaded")
}
