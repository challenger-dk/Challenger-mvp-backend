package routes

import (
	"server/controllers"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/users", func(r chi.Router) {
		r.Get("/", controllers.GetUsers)
	})
}
