# Challenger Backend

This is the Go backend service for the Challenger application. It's a monolithic API server responsible for handling all business logic, including users, challenges, teams, and invitations.

## Technology Stack

* **Language:** Go (Golang)
* **Web Framework:** Chi (v5)
* **Database:** PostgreSQL
* **ORM:** GORM
* **Authentication:** JWT (JSON Web Tokens)
* **Validation:** `go-playground/validator` (with custom rules)
* **Containerization:** Docker & Docker Compose

---

## Project Structure

The project code is located entirely within the `src/` directory, following a layered architecture:

* `/src`: The root of the Go module.
    * `main.go`: The main application entry point.
    * `go.mod`: Project dependencies.
    * `/config`: Handles environment variables (`env.go`), database connections (`db.go`), and sport seeding/caching (`sport.go`).
    * `/controllers`: HTTP handlers responsible for parsing requests, calling services, and writing HTTP responses.
    * `/services`: Contains all business logic (e.g., creating a user, accepting an invitation).
    * `/models`: Defines the GORM database models (structs) for tables like `User`, `Team`, `Challenge`, etc.
    * `/dto`: Data Transfer Objects used for request validation (`ChallengeCreateDto`) and response formatting (`ChallengeResponseDto`).
    * `/routes`: Defines all API endpoints using the Chi router and groups them by resource.
    * `/middleware`: Contains custom HTTP middleware like `AuthMiddleware` and `JsonContentType`.
    * `/tests`: Contains unit tests and integration tests.

---

## Core Features

* **Authentication:** Full JWT-based authentication with registration and login endpoints.
* **User Management:** Get users, get current user, update user profile (including favorite sports).
* **Teams:** Full CRUD operations for teams.
* **Challenges:** Full CRUD operations for challenges.
* **Invitations:** A complete system for sending, accepting, and declining invitations to resources (currently "team").
* **Dynamic Validation:** Includes custom validation rules, such as `is-valid-sport`, which dynamically checks against a list of sports loaded from the database at startup.

---

## Getting Started

You can run the project in two ways: with Docker (recommended) or locally.

### 1. Configuration

First, you must create a configuration file.

1.  Navigate to the `src/` directory.
2.  Copy the `.env.sample` file to a new file named `.env`.
3.  Fill in the values in `.env`. All fields are required.
