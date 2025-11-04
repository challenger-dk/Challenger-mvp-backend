package main

import (
	"log"
	"net/http"

	"server/config"
	"server/models"
	"server/services"

	"server/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Team{}, &models.Challenge{}, &models.Sport{})

	// Seed allowed sports
	if err := services.SeedSports(); err != nil {
		log.Fatal("Failed to seed sports:", err)
	}

	r := chi.NewRouter()

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://example.com", "http://localhost:3000"}, // Specify your frontend URLs
		AllowedOrigins:   []string{"*"}, // Allow all origins (use specific origins in production)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false, // Set to true if you need credentials, but can't use "*" with credentials
		MaxAge:           300,
	}))

	routes.RegisterRoutes(r)

	http.ListenAndServe(":8080", r)
}
