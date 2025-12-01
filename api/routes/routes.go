package routes

import (
	"server/api/controllers"
	"server/api/middleware"

	"github.com/go-chi/chi/v5"
	ext_middleware "github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(r chi.Router) {
	registerGenericMiddleware(r)

	// Public auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", controllers.Register)
		r.Post("/login", controllers.Login)
	})

	r.Get("/sports", controllers.GetSports)

	r.Route("/users", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/", controllers.GetUsers)
		r.Get("/me", controllers.GetCurrentUser)
		r.Get("/{id}", controllers.GetUserByID)
		r.Get("/settings", controllers.GetCurrentUserSettings)
		r.Get("/{id}/in-common", controllers.GetInCommonStats)
		r.Post("/block/{id}", controllers.BlockUser)
		r.Post("/unblock/{id}", controllers.UnblockUser)
		r.Put("/", controllers.UpdateUser)
		r.Put("/settings", controllers.UpdateUserSettings)
		r.Delete("/{id}", controllers.DeleteUser)
		r.Delete("/{id}/remove", controllers.RemoveFriend)
	})

	r.Route("/challenges", func(r chi.Router) {
		r.Get("/", controllers.GetChallenges)
		r.Get("/{id}", controllers.GetChallenge)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Post("/", controllers.CreateChallenge)
			r.Put("/{id}", controllers.UpdateChallenge)
			r.Post("/{id}/join", controllers.JoinChallenge)
			r.Post("/{id}/leave", controllers.LeaveChallenge)
			r.Delete("/{id}", controllers.DeleteChallenge)
		})
	})

	r.Route("/teams", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/{id}", controllers.GetTeam)
		r.Get("/", controllers.GetTeams)
		r.Get("/user/{id}", controllers.GetTeamsByUserId)
		r.Get("/me", controllers.GetCurrentUserTeams)

		r.Post("/", controllers.CreateTeam)
		//r.Post("/{id}/user", controllers.AddUserToTeam)

		r.Put("/{id}", controllers.UpdateTeam)

		r.Delete("/{id}", controllers.DeleteTeam)
		r.Delete("/{id}/user/{rmvUserId}", controllers.RemoveUserFromTeam)
		r.Delete("/{id}/leave", controllers.LeaveTeam)
	})

	r.Route("/invitations", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/user/{id}", controllers.GetInvitationsByUserId)
		r.Get("/me", controllers.GetCurrentUserInvitations)
		r.Post("/", controllers.SendInvitation)
		r.Post("/{id}/accept", controllers.AcceptInvitation)
		r.Post("/{id}/decline", controllers.DeclineInvitation)
	})

	r.Route("/notifications", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/", controllers.GetMyNotifications)
		r.Put("/read-all", controllers.MarkAllRead)
		r.Put("/{id}/read", controllers.MarkRead)
	})
}

func registerGenericMiddleware(r chi.Router) {
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.SlogMiddleware)
	r.Use(ext_middleware.RequestID) // Usefull for logging and tracing
	r.Use(ext_middleware.Recoverer)
	r.Use(ext_middleware.Heartbeat("/health"))
	r.Use(middleware.JsonContentType) // Sets Content-Type to json for all requests
}
