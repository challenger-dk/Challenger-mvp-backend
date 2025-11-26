.PHONY: seed help

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

seed: ## Seed the local database with test data
	@echo "🌱 Seeding database..."
	DB_HOST=localhost DB_PORT=5432 DB_USER=user DB_PASSWORD=password DB_NAME=challenger go run cmd/seed/main.go

# Install Atlas CLI (Mac/Linux)
atlas-install:
	curl -sSf https://atlasgo.sh | sh

# Generate a new migration file based on your Go structs
# Usage: make migrate-diff name=add_user_fields
migrate-diff:
	atlas migrate diff $(name) --env local

# Apply pending migrations to the database
migrate-apply:
	atlas migrate apply --env local --url "postgres://user:password@localhost:5432/challenger?sslmode=disable"