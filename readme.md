# Challenger Backend

This is the Go backend service for the Challenger application. It is structured as a modular monolith, separating the HTTP API and the WebSocket Chat service while sharing common logic and database models.

## Technology Stack

* **Language:** Go (1.24+)
* **Web Framework:** Chi (v5)
* **Database:** PostgreSQL 16 (with PostGIS extension)
* **ORM:** GORM
* **Migrations:** Atlas
* **Real-time:** Gorilla WebSocket
* **Authentication:** JWT (JSON Web Tokens)
* **Validation:** `go-playground/validator`
* **Containerization:** Docker & Docker Compose

---

## Services Architecture

The backend is architected as two distinct services that run concurrently. They share the same PostgreSQL database and common code modules (Models, DTOs, Config) located in the `/common` directory.

### 1. Main API Service (`/api`)
This is the primary service handling the core business logic of the application.
* **Responsibilities:** Manages User Authentication (JWT), User Profiles, Teams, Challenges, Sports, and the Invitation system.
* **Protocol:** RESTful HTTP (JSON).
* **Port:** Exposed on port `8000` (mapped via Docker).
* **Key Endpoints:** `/auth`, `/users`, `/teams`, `/challenges`, `/invitations`.

### 2. Chat Service (`/chat`)
A dedicated service designed to handle persistent connections and real-time messaging, keeping the main API lightweight.
* **Responsibilities:** Manages WebSocket connections for real-time communication in Team channels and Direct Messages. It also provides an HTTP endpoint to fetch message history.
* **Protocol:** WebSocket (for live events) & HTTP (for history).
* **Port:** Exposed on port `8002` (Development) or `8081` (Production).

---

## Project Structure

The project code is organized into the following directories:

* `/api`: Contains the REST API application code (Controllers, Services, Routes, Middleware).
* `/chat`: Contains the WebSocket Chat service code (Hub, Clients).
* `/common`: Shared code used by both the API and Chat services (Config, Models, DTOs).
* `cmd`: Command line utilities (`/atlas`, `/seed`).
* `Database`: Contains SQL schema initialization and migration files.

---

## 🚀 Getting Started

### Prerequisites
- Go 1.24 or higher
- Docker and Docker Compose
- Make (for running Makefile commands)
- Atlas CLI (for migrations) - Install with: `make atlas-install`

### 1. Configuration
Copy the sample environment file and configure it:
```bash
cp .env.sample .env
```

Edit `.env` and configure your environment variables:
- **Database credentials** (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
- **JWT secret** (JWT_SECRET)
- **Auto-migration** (AUTO_MIGRATE) - See [Database Migrations](#database-migrations) section

### 2. Start the Database
Start PostgreSQL with Docker Compose:
```bash
docker-compose up -d postgres
```

### 3. Run Database Migrations
Apply the database schema migrations:
```bash
make migrate-up
```

See the [Database Migrations](#database-migrations) section for more details.

### 4. Seed the Database (Optional)
Populate the database with test data:
```bash
make seed
```

### 5. Start the Services
Start both the API and Chat services:
```bash
# Terminal 1 - API Service
cd api
go run main.go

# Terminal 2 - Chat Service
cd chat
go run main.go
```

Or use Docker Compose to run everything:
```bash
docker-compose up
```

---

## 📊 Database Migrations

This project uses [Atlas](https://atlasgo.io/) for database schema migrations. Atlas provides a modern, declarative approach to managing database schemas with powerful features like automatic migration generation, rollback support, and pre-migration checks.

### Migration Philosophy

**Migrations are NOT run automatically by default.** This gives you full control over when schema changes are applied to your database.

- **Development**: You can enable auto-migration by setting `AUTO_MIGRATE=true` in your `.env` file for convenience
- **Production**: Always run migrations manually using the commands below

### Available Migration Commands

```bash
# Check migration status (what's pending, what's applied)
make migrate-status

# Preview migrations without applying them (dry-run)
make migrate-dry-run

# Apply pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Rollback ALL migrations (WARNING: destructive!)
make migrate-down-all

# Generate a new migration based on model changes
make migrate-diff name=add_user_field

# Validate migration files
make migrate-validate

# Rehash migration files (after manual edits)
make migrate-hash

# Lint migration files for potential issues
make migrate-lint
```

### Local Development Workflow

#### 1. Making Schema Changes

When you modify GORM models in `/common/models`, you need to generate a migration:

```bash
# 1. Make changes to your models (e.g., add a field to User struct)
# 2. Generate a migration
make migrate-diff name=add_user_email_verified

# 3. Review the generated migration file in /Database/migrations/
# 4. Apply the migration
make migrate-up
```

#### 2. Checking Migration Status

Before applying migrations, always check what's pending:

```bash
make migrate-status
```

Example output:
```
Migration Status: PENDING
  Current Version: 20240320120000
  Next Version:    20240320130000
  Pending:         1 migration
```

#### 3. Previewing Migrations (Dry Run)

Preview what will be applied without actually changing the database:

```bash
make migrate-dry-run
```

This shows you the exact SQL statements that will be executed.

#### 4. Applying Migrations

Apply all pending migrations:

```bash
make migrate-up
```

Atlas will:
- ✅ Check the current database state
- ✅ Apply pending migrations in order
- ✅ Run pre-migration checks to prevent data loss
- ✅ Update the migration history table

#### 5. Rolling Back Migrations

If you need to revert a migration (e.g., during local development):

```bash
# Rollback the last migration
make migrate-down
```

**How Atlas Rollback Works:**

Unlike traditional "down migrations" that use pre-written down files, Atlas computes the rollback plan based on the **current database state**:

- ✅ Works even if the migration was partially applied
- ✅ Runs pre-migration checks to warn about data deletion
- ✅ Uses transactions when possible (PostgreSQL)
- ✅ Handles non-transactional databases safely (MySQL)

**Example:**
```bash
# Check what will be rolled back
make migrate-dry-run

# Rollback the last migration
make migrate-down

# After rollback, you can delete the migration file
atlas migrate rm 20240320130000

# Generate a new migration with fixes
make migrate-diff name=add_user_email_verified_fixed
```

### Common Scenarios

#### Scenario 1: Fresh Database Setup
```bash
# Start the database
docker-compose up -d postgres

# Apply all migrations
make migrate-up

# Seed with test data
make seed
```

#### Scenario 2: Experimenting with Schema Changes
```bash
# Make changes to models
# Generate migration
make migrate-diff name=experiment_user_fields

# Apply it
make migrate-up

# Test your changes...

# Oops, need to revert
make migrate-down

# Delete the migration file
atlas migrate rm <version>

# Try again with different approach
make migrate-diff name=experiment_user_fields_v2
make migrate-up
```

#### Scenario 3: Pulling Latest Changes
```bash
# Pull latest code (includes new migrations)
git pull

# Check what migrations are pending
make migrate-status

# Preview the changes
make migrate-dry-run

# Apply them
make migrate-up
```

#### Scenario 4: Resetting Your Local Database
```bash
# Option 1: Rollback all migrations
make migrate-down-all

# Option 2: Drop and recreate the database
docker-compose down -v
docker-compose up -d postgres
make migrate-up
make seed
```

### Migration Best Practices

1. **Always review generated migrations** before applying them
2. **Use descriptive names** when generating migrations: `make migrate-diff name=add_user_email_verification`
3. **Check migration status** before and after applying: `make migrate-status`
4. **Use dry-run** to preview changes: `make migrate-dry-run`
5. **Make migrations backwards compatible** when possible (additive changes)
6. **Never edit applied migrations** - create a new migration instead
7. **Commit migration files** to version control immediately after generating them

### Troubleshooting

#### Migration Hash Mismatch
If you see a hash mismatch error after manually editing a migration:
```bash
make migrate-hash
```

#### Migration Failed Midway
Atlas tracks which statements were applied. Simply re-run:
```bash
make migrate-up
```

#### Want to Start Fresh
```bash
# Drop the database and start over
docker-compose down -v
docker-compose up -d postgres
make migrate-up
make seed
```

---

## 🔍 Code Quality & Linting

### Sanitize Validation Linter

All string fields in input DTOs (Create/Update DTOs) must have the `sanitize` validation tag to prevent XSS attacks and SQL injection.

**Run the linter:**
```bash
make lint-sanitize
```

**What it checks:**
- ✅ All `string` fields in `*CreateDto` structs have `validate:"sanitize,..."`
- ✅ All `string` fields in `*UpdateDto` structs have `validate:"sanitize,..."`
- ✅ All `string` fields in `Login` struct have `validate:"sanitize,..."`
- ⏭️ Skips response DTOs (they contain already-sanitized data from the database)

**Example:**
```go
// ✅ Good
type TeamCreateDto struct {
    Name string `json:"name" validate:"sanitize,required,min=3"`
}

// ❌ Bad - missing sanitize
type TeamCreateDto struct {
    Name string `json:"name" validate:"required,min=3"`
}
```

**Run before committing:**
```bash
make lint-sanitize
```

**Automated Checks:**
- ✅ Runs automatically in GitHub Actions on every push/PR
- ✅ Runs before deployment to prevent insecure code from being deployed
- ❌ Deployment will fail if any string fields are missing `sanitize` tag

---

## 🧪 Testing

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./api/handlers
```

---

## 📁 Key Directories

- `/api` - REST API service
  - `/handlers` - HTTP request handlers
  - `/middleware` - Authentication, CORS, logging
  - `/routes` - Route definitions
  - `/services` - Business logic layer

- `/chat` - WebSocket chat service
  - `/handlers` - WebSocket and HTTP handlers
  - `/hub` - WebSocket connection manager
  - `/routes` - Route definitions

- `/common` - Shared code
  - `/config` - Configuration and database connection
  - `/models` - GORM database models
  - `/dto` - Data Transfer Objects

- `/Database/migrations` - Atlas migration files

- `/cmd` - Command-line utilities
  - `/seed` - Database seeding tool

---

## 🔧 Development Tips

### Auto-Migration for Development
For local development convenience, you can enable auto-migration:

```bash
# In your .env file
AUTO_MIGRATE=true
```

With this setting, migrations will run automatically when you start the API service. This is useful for rapid development but **should never be used in production**.

### Viewing Logs
Both services use structured logging. Key log messages include:
- 🔄 Auto-migration status
- 📊 Database connection status
- 🚀 Service startup and port binding
- ⚠️ Errors and warnings

### Database Access
Connect to the PostgreSQL database:
```bash
# Using Docker
docker exec -it challenger-postgres psql -U user -d challenger

# Using psql directly
psql -h localhost -p 5432 -U user -d challenger
```

---

## 🚢 Production Deployment

### Google Cloud Migrations (Cloud Run Jobs)

**One command to run migrations:**

```bash
# Staging
gcloud run jobs execute challenger-migrate-staging --region=us-central1

# Production
gcloud run jobs execute challenger-migrate-production --region=us-central1
```

See **[GOOGLE_CLOUD_MIGRATIONS.md](GOOGLE_CLOUD_MIGRATIONS.md)** for setup instructions.

### Migration Strategy

**Never use auto-migration in production.** Instead:

1. **Run migrations manually** before deploying new code
2. **Verify migration status**
3. **Deploy the application** after migrations are successful

### Environment Variables

Ensure these are set in production:
- `AUTO_MIGRATE=false` (or unset - this is the default)
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `JWT_SECRET` (use a strong, random secret)
- `ENVIRONMENT=production`

---

## 📚 Additional Resources

- [Atlas Documentation](https://atlasgo.io/docs)
- [GORM Documentation](https://gorm.io/docs/)
- [Chi Router Documentation](https://go-chi.io/)
- [Gorilla WebSocket Documentation](https://pkg.go.dev/github.com/gorilla/websocket)

---
