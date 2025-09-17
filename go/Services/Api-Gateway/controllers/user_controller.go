package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func routes(uc *UserController) http.Handler {
	r := chi.NewRouter()

	r.Get("")

}

// Handlers

func get_users(uc *UserController) (w http.ResponseWriter, r *http.Request) {

}
