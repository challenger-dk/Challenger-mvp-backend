package integration

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"

	"server/common/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var setupOnce sync.Once

// TestMain acts as the entry point for all tests in this package.
func TestMain(m *testing.M) {
	setupTestDB()
	code := m.Run()
	os.Exit(code)
}

func setupTestDB() {
	setupOnce.Do(func() {
		// Connection string for the Docker "postgres-test" container
		// Port 5433 matches the docker-compose.yml test service
		dsn := "host=localhost user=test_user password=test_password dbname=challenger_test port=5433 sslmode=disable"

		var err error
		// We explicitly assign to the global variable in the config package
		config.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Fatalf("❌ Failed to connect to test database: %v", err)
		}

		// Get raw SQL connection for Atlas migrations
		sqlDB, err := config.DB.DB()
		if err != nil {
			log.Fatalf("❌ Failed to get database connection: %v", err)
		}

		// Drop all schemas except system ones to start fresh
		// This ensures a clean state for migrations
		dropSchemas := []string{
			"DROP SCHEMA IF EXISTS atlas_schema_revisions CASCADE", // Drop Atlas migration tracking
			"DROP SCHEMA IF EXISTS public CASCADE",
			"DROP SCHEMA IF EXISTS tiger CASCADE",
			"DROP SCHEMA IF EXISTS tiger_data CASCADE",
			"DROP SCHEMA IF EXISTS topology CASCADE",
			"CREATE SCHEMA public",
		}
		for _, stmt := range dropSchemas {
			if err := config.DB.Exec(stmt).Error; err != nil {
				log.Fatalf("❌ Failed to execute: %s - %v", stmt, err)
			}
		}

		// Run Atlas migrations
		dbURL := "postgres://test_user:test_password@localhost:5433/challenger_test?sslmode=disable"
		if err := config.RunAtlasMigrationsWithConnection(sqlDB, dbURL); err != nil {
			log.Fatalf("❌ Failed to run Atlas migrations: %v", err)
		}

		// Seed basic data needed for tests
		if err := config.SeedSports(); err != nil {
			log.Fatalf("❌ Failed to seed sports: %v", err)
		}
	})
}

// setupTest prepares the DB for a specific test case (cleans it).
func setupTest(_ *testing.T) func() {
	// Ensure DB is set up even if TestMain didn't run (defensive coding)
	if config.DB == nil {
		setupTestDB()
	}

	// Clean before running the test to ensure a clean slate
	clearDB()

	// Return a cleanup function
	return func() {
		// Optional: Clear after as well, or leave data for inspection on failure
		// clearDB()
	}
}

func clearDB() {
	if config.DB == nil {
		panic("🔥 config.DB is nil in clearDB! Database connection failed or was not initialized.")
	}

	// Truncate tables in specific order to handle foreign keys
	tables := []string{
		"reports",
		"messages",
		"notifications",
		"invitations",
		"team_sports",
		"user_favorite_sports",
		"team_members",
		"user_friends",
		"challenge_teams",
		"user_challenges",
		"challenges",
		"teams",
		"user_settings",
		"users",
		"locations",
	}

	for _, table := range tables {
		if err := config.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error; err != nil {
			log.Printf("⚠️ Failed to truncate table %s: %v", table, err)
		}
	}
}
