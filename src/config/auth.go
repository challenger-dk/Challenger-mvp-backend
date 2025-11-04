package config

import (
	"os"
)

var (
	JWTSecret = getEnv("JWT_SECRET", "Lfi0rO+Qq2w2r6YiNlqngPOgr/gAYNu81k2b6SwqFM0=")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

