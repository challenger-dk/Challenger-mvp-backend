package routes

import (
	"server/controllers"
	"server/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	// Public auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", controllers.Register)
		r.Post("/login", controllers.Login)
	})

	r.Get("/sports", controllers.GetSports)

	// Protected routes
	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/", controllers.GetUsers)
		r.Get("/me", controllers.GetCurrentUser)
		r.Get("/{id}", controllers.GetUserByID)
		r.Put("/{id}", controllers.UpdateUser)
		r.Delete("/{id}", controllers.DeleteUser)
	})

	r.Route("/challenges", func(r chi.Router) {
		r.Get("/", controllers.GetChallenges)
		r.Get("/{id}", controllers.GetChallenge)

		r.Post("/", controllers.CreateChallenge)

		r.Put("/{id}", controllers.UpdateChallenge)

		r.Delete("/{id}", controllers.DeleteChallenge)
	})

	r.Route("/teams", func(r chi.Router) {
		r.Get("/{id}", controllers.GetTeam)
		r.Get("/", controllers.GetTeams)
		r.Post("/", controllers.CreateTeam)
		r.Put("/{id}", controllers.UpdateTeam)
		r.Delete("/{id}", controllers.DeleteTeam)
	})
}
