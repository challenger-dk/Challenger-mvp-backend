.PHONY: help seed migrate-diff migrate-up migrate-down migrate-down-all migrate-status migrate-hash migrate-validate migrate-dry-run atlas-install lint-sanitize test test-unit test-integration test-coverage test-coverage-unit test-coverage-integration test-coverage-summary

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

# ============================================================================
# Testing
# ============================================================================

test: ## Run all tests without coverage
	@echo "🧪 Running all tests..."
	@go test -v ./tests/... ./common/... ./api/... ./chat/...

test-unit: ## Run unit tests only
	@echo "🧪 Running unit tests..."
	@go test -v ./tests/unit/...

test-integration: ## Run integration tests only
	@echo "🧪 Running integration tests..."
	@go test -v ./tests/integration/...

# ============================================================================
# Test Coverage
# ============================================================================

test-coverage: ## Generate coverage profile and show summary
	@echo "🧪 Running all tests with coverage..."
	@go test -coverprofile=coverage.out -covermode=atomic \
		-coverpkg=./common/...,./api/...,./chat/... \
		./tests/... > /dev/null 2>&1
	@echo ""
	@echo "📊 Coverage Summary:"
	@go tool cover -func=coverage.out | tail -1
	@echo ""
	@echo "✅ Coverage profile generated: coverage.out"

test-coverage-unit: ## Generate coverage profile for unit tests only
	@echo "🧪 Running unit tests with coverage..."
	@go test -coverprofile=coverage-unit.out -covermode=atomic \
		-coverpkg=./common/...,./api/...,./chat/... \
		./tests/unit/... > /dev/null 2>&1
	@echo ""
	@echo "📊 Unit Test Coverage:"
	@go tool cover -func=coverage-unit.out | tail -1
	@echo "✅ Coverage profile generated: coverage-unit.out"

test-coverage-integration: ## Generate coverage profile for integration tests only
	@echo "🧪 Running integration tests with coverage..."
	@go test -coverprofile=coverage-integration.out -covermode=atomic \
		-coverpkg=./common/...,./api/...,./chat/... \
		./tests/integration/... > /dev/null 2>&1
	@echo ""
	@echo "📊 Integration Test Coverage:"
	@go tool cover -func=coverage-integration.out | tail -1
	@echo "✅ Coverage profile generated: coverage-integration.out"

test-coverage-summary: ## Show detailed coverage breakdown by package
	@if [ ! -f coverage.out ]; then \
		echo "❌ coverage.out not found. Run 'make test-coverage' first."; \
		exit 1; \
	fi
	@echo "📊 Coverage Breakdown by Package:"
	@go tool cover -func=coverage.out | grep "^server/" | grep -E "common|api|chat" | sort -t$$'\t' -k3 -rn
	@echo ""
	@echo "📈 Overall Coverage:"
	@go tool cover -func=coverage.out | tail -1