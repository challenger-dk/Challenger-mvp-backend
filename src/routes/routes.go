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

	// Protected routes
	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/", controllers.GetUsers)
		r.Get("/me", controllers.GetCurrentUser)
	})
}
