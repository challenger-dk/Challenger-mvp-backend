package main

import (
	"log"
	"net/http"
	"server/common/config"
	"server/common/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	config.LoadConfig()
	config.ConnectDatabase()

	// Ensure tables exist
	err := config.DB.AutoMigrate(&models.Message{}, &models.User{}, &models.Team{}, &models.Chat{})
	if err != nil {
		log.Fatal(err)
	}

	hub := newHub()
	go hub.run()

	mux := http.NewServeMux()

	// WebSocket Endpoint
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// Message History API Endpoint
	mux.HandleFunc("/api/messages", getMessages)

	server := &http.Server{
		Addr:         ":8002",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Chat Service started on :8002")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("ListenAndServe: ", err)
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
	// Preload Teams, Chats, and BlockedUsers
	if err := config.DB.Preload("Teams").Preload("Chats").Preload("BlockedUsers").First(&user, claims.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	allowedTeams := make(map[uint]bool)
	for _, team := range user.Teams {
		allowedTeams[team.ID] = true
	}

	allowedChats := make(map[uint]bool)
	for _, chat := range user.Chats {
		allowedChats[chat.ID] = true
	}

	blockedUsers := make(map[uint]bool)
	for _, blocked := range user.BlockedUsers {
		blockedUsers[blocked.ID] = true
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:            hub,
		conn:           conn,
		send:           make(chan []byte, 256),
		userID:         claims.UserID,
		teamIDs:        allowedTeams,
		chatIDs:        allowedChats,
		blockedUserIDs: blockedUsers,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
