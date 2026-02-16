package main

import (
	"log/slog"
	"net/http"
	"os"
	"server/chat/handlers"
	"server/common/config"
	"server/common/logger"
	commonMiddleware "server/common/middleware"
	"server/common/models"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// 1. Initialize Global Logger first
	logger.InitLogger()
	slog.Info("💬 Chat Service starting...")

	config.LoadConfig()
	config.ConnectDatabase()

	// Note: Database migrations and PostGIS extension are handled by the API service
	// The chat service only needs to connect to the database
	slog.Info("✅ Database connected")

	hub := newHub()
	go hub.run()

	// Setup Chi router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// WebSocket Endpoint (no auth middleware, uses query param token)
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// Legacy Message History API Endpoint
	r.Get("/api/messages", getMessages)

	// New Conversation API Routes (with auth middleware)
	r.Route("/api/conversations", func(r chi.Router) {
		r.Use(commonMiddleware.AuthMiddleware)
		r.Use(commonMiddleware.EulaMiddleware)

		r.Post("/direct", handlers.CreateDirectConversation)
		r.Post("/group", handlers.CreateGroupConversation)
		r.Get("/", handlers.ListConversations)
		r.Get("/{id}", handlers.GetConversation)
		r.Get("/{id}/messages", handlers.GetConversationMessages)
		r.Post("/{id}/messages", handlers.SendMessage)
		r.Post("/{id}/read", handlers.MarkConversationRead)
		r.Get("/team/{teamId}", handlers.GetTeamConversation)
		r.Get("/challenge/{challengeId}", handlers.GetChallengeConversation)
	})

	// Internal endpoints for team/challenge sync (no auth for internal service calls)
	r.Post("/internal/teams/{teamId}/sync", handlers.SyncTeamMembers)
	r.Post("/internal/challenges/{challengeId}/sync", handlers.SyncChallengeMembers)

	// Port from Cloud Run environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	slog.Info("Starting server on :" + port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	claims := &struct {
		UserID uint `json:"user_id"`
		jwt.RegisteredClaims
	}{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := config.DB.Preload("Teams.Team").Preload("BlockedUsers").First(&user, claims.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	allowedTeams := make(map[uint]bool)
	for _, team := range user.Teams {
		allowedTeams[team.TeamID] = true
	}

	blockedUsers := make(map[uint]bool)
	for _, blocked := range user.BlockedUsers {
		blockedUsers[blocked.ID] = true
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", "error", err)
		return
	}

	client := &Client{
		hub:            hub,
		conn:           conn,
		send:           make(chan []byte, 256),
		userID:         claims.UserID,
		teamIDs:        allowedTeams,
		blockedUserIDs: blockedUsers,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
