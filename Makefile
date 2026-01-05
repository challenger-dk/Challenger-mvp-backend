.PHONY: help seed migrate-diff migrate-up migrate-down migrate-down-all migrate-status migrate-hash migrate-validate migrate-dry-run atlas-install lint-sanitize

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ============================================================================
# Database Seeding
# ============================================================================

seed: ## Seed the local database with test data
	@echo "🌱 Seeding database..."
	DB_HOST=localhost DB_PORT=5432 DB_USER=user DB_PASSWORD=password DB_NAME=challenger go run cmd/seed/main.go

# ============================================================================
# Atlas Migration Commands
# ============================================================================

atlas-install: ## Install Atlas CLI (Mac/Linux/WSL)
	@echo "📦 Installing Atlas CLI..."
	curl -sSf https://atlasgo.sh | sh

migrate-diff: ## Generate a new migration file based on your Go structs (Usage: make migrate-diff name=add_user_fields)
	@echo "🔍 Generating migration diff..."
	atlas migrate diff $(name) --env local

migrate-up: ## Apply pending migrations to the local database
	@echo "⬆️  Applying migrations..."
	atlas migrate apply --env local --url "postgres://user:password@localhost:5432/challenger?sslmode=disable"

migrate-down: ## Rollback the last migration (Usage: make migrate-down)
	@echo "⬇️  Rolling back last migration..."
	atlas migrate down --env local --url "postgres://user:password@localhost:5432/challenger?sslmode=disable"

migrate-down-all: ## Rollback all migrations (WARNING: This will drop all tables!)
	@echo "⚠️  WARNING: This will rollback ALL migrations!"
	@echo "Press Ctrl+C to cancel, or wait 5 seconds to continue..."
	@sleep 5
	atlas migrate down --env local --url "postgres://user:password@localhost:5432/challenger?sslmode=disable" --to-version 0

migrate-dry-run: ## Preview what migrations would be applied without actually applying them
	@echo "🔍 Dry run - previewing migrations..."
	atlas migrate apply --env local --url "postgres://user:password@localhost:5432/challenger?sslmode=disable" --dry-run

migrate-status: ## Check migration status for local database
	@echo "📊 Checking migration status..."
	atlas migrate status --env local --url "postgres://user:password@localhost:5432/challenger?sslmode=disable"

migrate-hash: ## Rehash migration files (run after manually editing migrations)
	@echo "🔐 Rehashing migration files..."
	atlas migrate hash --env local

migrate-validate: ## Validate migration files
	@echo "✅ Validating migration files..."
	atlas migrate validate --env local

migrate-lint: ## Lint migration files for potential issues
	@echo "🔍 Linting migration files..."
	atlas migrate lint --env local --latest 1

# ============================================================================
# Code Quality & Linting
# ============================================================================

lint-sanitize: ## Check that all string fields in DTOs have sanitize validation tags
	@echo "🔍 Checking sanitize validation tags..."
	@go run tools/validate-sanitize/main.go