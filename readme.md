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

### 1. Configuration
Copy the sample environment file and configure it.
```bash
cp .env.sample .env
