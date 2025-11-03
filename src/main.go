package main

import (
	"net/http"

	"server/config"
	"server/models"

	"server/controllers"

	"github.com/go-chi/chi/v5"
)

func main() {

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Team{}, &models.Challenge{})
	r := chi.NewRouter()

	r.Get("/", controllers.GetUsers)

	http.ListenAndServe(":8080", r)
}
