package main

import (
	"net/http"

	"server/config"
	"server/models"

	"server/routes"

	"github.com/go-chi/chi/v5"
)

func main() {

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Team{}, &models.Challenge{})
	r := chi.NewRouter()
	routes.RegisterRoutes(r)

	http.ListenAndServe(":8080", r)
}
