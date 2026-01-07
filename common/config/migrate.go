package config

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"

	"ariga.io/atlas-go-sdk/atlasexec"
)

// getMigrationsDir returns the file URL to the migrations directory
func getMigrationsDir() (string, error) {
	// Try to find the migrations directory by walking up from the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Walk up the directory tree to find the Database/migrations directory
	dir := cwd
	for {
		migrationsPath := filepath.Join(dir, "Database", "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			// Convert to forward slashes for file URL
			migrationsPath = filepath.ToSlash(migrationsPath)
			return fmt.Sprintf("file://%s", migrationsPath), nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the root without finding migrations
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find Database/migrations directory")
}

// RunAtlasMigrations applies pending Atlas migrations to the database
func RunAtlasMigrations() error {
	slog.Info("🔄 Running Atlas migrations...")

	// Get database connection string
	// Note: In production, Atlas uses database locks to prevent concurrent migrations
	// Only the first instance to acquire the lock will run migrations
	// Other instances will wait and then skip if migrations are already applied

	var dsn string

	// Check if using Unix socket (Cloud SQL) or TCP connection
	if len(AppConfig.DBHost) > 0 && AppConfig.DBHost[0] == '/' {
		// Unix socket connection (Cloud SQL)
		// Format: postgres://user:password@/dbname?host=/cloudsql/instance-connection-name
		dsn = fmt.Sprintf("postgres://%s:%s@/%s?host=%s&sslmode=disable",
			url.QueryEscape(AppConfig.DBUser),
			url.QueryEscape(AppConfig.DBPassword),
			AppConfig.DBName,
			url.QueryEscape(AppConfig.DBHost),
		)
		slog.Info("🔍 Using Unix socket connection", "socket", AppConfig.DBHost)
	} else {
		// TCP connection
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			url.QueryEscape(AppConfig.DBUser),
			url.QueryEscape(AppConfig.DBPassword),
			AppConfig.DBHost,
			AppConfig.DBPort,
			AppConfig.DBName,
		)
		slog.Info("🔍 Using TCP connection", "host", AppConfig.DBHost, "port", AppConfig.DBPort)
	}

	// Debug: Log connection details (without password)
	slog.Info("🔍 Migration DSN details",
		"user", AppConfig.DBUser,
		"database", AppConfig.DBName,
	)

	// Ensure extensions exist before running migrations
	if err := ensureExtensions(); err != nil {
		return fmt.Errorf("failed to ensure extensions: %w", err)
	}

	// Get migrations directory
	migrationsDirURL, err := getMigrationsDir()
	if err != nil {
		return fmt.Errorf("failed to find migrations directory: %w", err)
	}

	// Create Atlas client
	client, err := atlasexec.NewClient(".", "atlas")
	if err != nil {
		return fmt.Errorf("failed to create atlas client: %w", err)
	}

	ctx := context.Background()

	// Apply migrations
	result, err := client.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL:    dsn,
		DirURL: migrationsDirURL,
	})

	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if result.Error != "" {
		return fmt.Errorf("migration error: %s", result.Error)
	}

	slog.Info("✅ Atlas migrations applied successfully",
		"target", result.Target,
		"current", result.Current,
		"applied", len(result.Applied),
	)

	return nil
}

// ensureExtensions creates required PostgreSQL extensions
func ensureExtensions() error {
	// Get raw SQL connection from GORM
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Create required extensions
	extensions := []string{"postgis", "pg_trgm"}
	for _, ext := range extensions {
		query := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s", ext)
		if _, err := sqlDB.Exec(query); err != nil {
			return fmt.Errorf("failed to create extension %s: %w", ext, err)
		}
		slog.Info("✅ Extension ensured", "extension", ext)
	}

	return nil
}

// RunAtlasMigrationsWithConnection applies migrations using an existing SQL connection
// This is useful for tests where we want to use a specific database connection
func RunAtlasMigrationsWithConnection(sqlDB *sql.DB, dbURL string) error {
	slog.Info("🔄 Running Atlas migrations with custom connection...")

	// Ensure extensions exist
	extensions := []string{"postgis", "pg_trgm"}
	for _, ext := range extensions {
		query := fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %s", ext)
		if _, err := sqlDB.Exec(query); err != nil {
			return fmt.Errorf("failed to create extension %s: %w", ext, err)
		}
	}

	// Get migrations directory
	migrationsDirURL, err := getMigrationsDir()
	if err != nil {
		return fmt.Errorf("failed to find migrations directory: %w", err)
	}

	// Create Atlas client
	client, err := atlasexec.NewClient(".", "atlas")
	if err != nil {
		return fmt.Errorf("failed to create atlas client: %w", err)
	}

	// Apply migrations
	ctx := context.Background()
	result, err := client.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL:    dbURL,
		DirURL: migrationsDirURL,
	})

	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if result.Error != "" {
		return fmt.Errorf("migration error: %s", result.Error)
	}

	slog.Info("✅ Atlas migrations applied successfully",
		"target", result.Target,
		"current", result.Current,
		"applied", len(result.Applied),
	)

	return nil
}
