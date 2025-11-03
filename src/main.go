package main

import (
	"net/http"

    "server/config"
    "server/models"
	
	"github.com/go-chi/chi/v5"
)

func main() {

	config.ConnectDatabase()
    if err := config.DB.AutoMigrate(&models.User{}); err != nil {
        // Fail fast if migration cannot run
        panic(err)
    }
	r := chi.NewRouter()

	http.ListenAndServe(":8080", r)
}
