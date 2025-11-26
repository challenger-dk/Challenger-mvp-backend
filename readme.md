# Challenger Backend

This is the Go backend service for the Challenger application. It is structured as a modular monolith, separating the HTTP API and the WebSocket Chat service while sharing common logic.

## Technology Stack

* **Language:** Go (Golang)
* **Web Framework:** Chi (v5)
* **Database:** PostgreSQL (with PostGIS)
* **ORM:** GORM
* **Authentication:** JWT (JSON Web Tokens)
* **Validation:** `go-playground/validator` (with custom rules)
* **Containerization:** Docker & Docker Compose

---

## Services Architecture

The backend is architected as two distinct services that run concurrently. They share the same PostgreSQL database and common code modules (Models, DTOs, Config) located in the `/common` directory.

### 1. Main API Service (`/api`)
This is the primary service handling the core business logic of the application.
* **Responsibilities:** Manages User Authentication (JWT), User Profiles, Teams, Challenges, Sports, and the Invitation system.
* **Protocol:** RESTful HTTP (JSON).
* **Port:** Exposed on port `8080`.
* **Key Endpoints:** `/auth`, `/users`, `/teams`, `/challenges`, `/invitations`.

### 2. Chat Service (`/chat`)
A dedicated service designed to handle persistent connections and real-time messaging, keeping the main API lightweight.
* **Responsibilities:** Manages WebSocket connections for real-time communication in Team channels and Direct Messages. It also provides an HTTP endpoint to fetch message history.
* **Protocol:** WebSocket (for live events) & HTTP (for history).
* **Port:** Exposed on port `8082` (Development) or `8081` (Production).

---

## Project Structure

The project code is organized into the following directories:

* `/api`: Contains the REST API application code.
    * `main.go`: Entry point for the API server.
    * `/controllers`: HTTP handlers for parsing requests and managing responses.
    * `/services`: Business logic for the API.
    * `/routes`: API route definitions.
    * `/middleware`: HTTP middleware (Auth, CORS, etc.).
    * `/tests`: Unit and integration tests.
* `/chat`: Contains the WebSocket Chat service code.
    * `main.go`: Entry point for the Chat server.
    * `hub.go` & `client.go`: WebSocket logic for handling real-time connections.
* `/common`: Shared code used by both the API and Chat services.
    * `/config`: Configuration (`env.go`), database connections (`db.go`), and setup.
    * `/models`: GORM database models (`User`, `Team`, `Challenge`, etc.).
    * `/dto`: Shared Data Transfer Objects.
    * `/appError`: Centralized error handling.
* `go.mod`: Project dependencies.
* `docker-compose.yml`: Local development environment setup.

---

## Core Features

* **Authentication:** Full JWT-based authentication with registration and login endpoints.
* **User Management:** Get users, get current user, update user profile (including favorite sports).
* **Teams:** Full CRUD operations for teams.
* **Challenges:** Full CRUD operations for challenges.
* **Invitations:** A complete system for sending, accepting, and declining invitations.
* **Real-time Chat:** WebSocket-based chat for teams and direct messages.
* **Dynamic Validation:** Includes custom validation rules, such as `is-valid-sport`.

---

## Getting Started

You can run the project in two ways: with Docker (recommended) or locally.

### 1. Configuration

First, you must create a configuration file.

1.  Navigate to the root directory.
2.  Copy the `.env.sample` file to a new file named `.env`.
3.  Fill in the values in `.env`. All fields are required.

### 2. Running with Docker

1.  Ensure Docker and Docker Compose are installed.
2.  Run the application:
    ```bash
    docker compose up --build
    ```
3.  The services will be available at:
    * **API:** `http://localhost:8080`
    * **Chat:** `http://localhost:8082` (WebSocket endpoint at `ws://localhost:8082/ws`)

4.  Run a clean seed to create dummy data for local envoriment
    ```bash
    make seed
    ```