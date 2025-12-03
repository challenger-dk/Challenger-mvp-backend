package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/common/config"
	"server/common/dto"
	"server/common/models"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

// Helper function to extract and validate claims from the JWT token
func authenticateRequest(r *http.Request) (*struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" || len(tokenString) < 7 || tokenString[:6] != "Bearer" {
		return nil, fmt.Errorf("missing or invalid Authorization header")
	}
	tokenString = tokenString[7:] // Remove "Bearer "

	claims := &struct {
		UserID uint `json:"user_id"`
		jwt.RegisteredClaims
	}{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

// getMessages handles fetching chat history.
func getMessages(w http.ResponseWriter, r *http.Request) {
	claims, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 1. Get Query Parameters
	teamIDStr := r.URL.Query().Get("team_id")
	chatIDStr := r.URL.Query().Get("chat_id")
	userID := claims.UserID

	limit := 50

	// 2. Build the DB Query
	var messages []models.Message
	query := config.DB.Preload("Sender").Order("created_at DESC").Limit(limit)

	if teamIDStr != "" {
		// --- Team Chat History ---
		teamID, err := strconv.ParseUint(teamIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid team_id format", http.StatusBadRequest)
			return
		}

		// Authorization Check
		var isMember bool
		config.DB.Model(&models.Team{}).
			Joins("JOIN user_teams ON user_teams.team_id = teams.id").
			Where("teams.id = ? AND user_teams.user_id = ?", teamID, userID).
			Select("count(*) > 0").
			Find(&isMember)

		if !isMember {
			http.Error(w, "You are not a member of this team", http.StatusForbidden)
			return
		}

		query = query.Where("team_id = ?", uint(teamID))

	} else if chatIDStr != "" {
		// --- Group/DM Chat History ---
		chatID, err := strconv.ParseUint(chatIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid chat_id format", http.StatusBadRequest)
			return
		}

		// Authorization Check
		var isMember bool
		config.DB.Model(&models.Chat{}).
			Joins("JOIN user_chats ON user_chats.chat_id = chats.id").
			Where("chats.id = ? AND user_chats.user_id = ?", chatID, userID).
			Select("count(*) > 0").
			Find(&isMember)

		if !isMember {
			http.Error(w, "You are not a participant of this chat", http.StatusForbidden)
			return
		}

		query = query.Where("chat_id = ?", uint(chatID))

	} else {
		http.Error(w, "Must provide either team_id or chat_id", http.StatusBadRequest)
		return
	}

	// 3. Execute the Query
	if err := query.Find(&messages).Error; err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}

	// 4. Transform to DTOs and Reverse Order
	responseDTOs := make([]dto.MessageResponseDto, len(messages))
	for i, msg := range messages {
		responseDTOs[len(messages)-1-i] = dto.ToMessageResponseDto(msg)
	}

	// 5. Send Response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseDTOs); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
