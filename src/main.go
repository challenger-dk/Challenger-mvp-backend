package main

import (
	"net/http"

    "server/config"
    "server/models"
	
	"github.com/go-chi/chi/v5"
)

func main() {

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Team{}, &models.Challenge{})
	r := chi.NewRouter()

	http.ListenAndServe(":8080", r)
}
