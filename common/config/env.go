package config

import (
	"log/slog"
	"os"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	// App Environment (development, production)
	AppEnv string `env:"APP_ENV" envDefault:"development"`

	// Database connection settings
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,required"`

	// JWT Secret
	JWTSecret string `env:"JWT_SECRET,required"`

	// JWT Token Expiration (in hours, default 30 days = 720 hours)
	JWTExpirationHours int `env:"JWT_EXPIRATION_HOURS" envDefault:"720"`

	// Cron Settings
	EnableCron bool `env:"ENABLE_CRON" envDefault:"true"`

	// Postmark API Key
	PostmarkAPIKey string `env:"POSTMARK_API_KEY,required"`

	// Postmark From Email
	PostmarkFromEmail string `env:"POSTMARK_FROM_EMAIL,required"`

	// Firebase Project ID (for OAuth token verification)
	FirebaseProjectID string `env:"FIREBASE_PROJECT_ID,required"`

	// Weather API Key
	WeatherAPIKey string `env:"WEATHER_API_KEY"`
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment")
	}

	if err := env.Parse(&AppConfig); err != nil {
		slog.Error("Failed to parse config", "error", err)
		os.Exit(1)
	}

	slog.Info("Configuration loaded", "env", AppConfig.AppEnv)
}
