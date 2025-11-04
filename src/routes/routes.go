package routes

import (
	"server/controllers"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {

	r.Route("/users", func(r chi.Router) {
		r.Get("/", controllers.GetUsers)
	})

	r.Route("/teams", func(r chi.Router) {
		r.Get("/{id}", controllers.GetTeam)
		r.Get("/", controllers.GetTeams)
		r.Post("/", controllers.CreateTeam)
		r.Put("/{id}", controllers.UpdateTeam)
		r.Delete("/{id}", controllers.DeleteTeam)
	})
}
