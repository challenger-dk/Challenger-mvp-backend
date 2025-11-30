package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/common/config"
	"server/common/dto"
	"server/common/models"
	"strconv"
	"math"
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

// getMessages handles fetching chat history for a group or 1-on-1 conversation.
func getMessages(w http.ResponseWriter, r *http.Request) {
	claims, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 1. Get Query Parameters
	teamIDStr := r.URL.Query().Get("team_id")
	recipientIDStr := r.URL.Query().Get("recipient_id")
	userID := claims.UserID

	// Default limit for messages
	limit := 50

	// 2. Build the DB Query
	var messages []models.Message
	query := config.DB.Preload("Sender").Order("created_at DESC").Limit(limit)

	if teamIDStr != "" {
		// --- Group Chat History ---
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

		// A real-world check would verify if 'userID' is actually a member of 'teamID'.
		// Assuming for now the user is authorized if the token is valid.
		query = query.Where("team_id = ?", uint(teamID))

	} else if recipientIDStr != "" {
		// --- 1-on-1 Chat History ---
		recipientID, err := strconv.ParseUint(recipientIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid recipient_id format", http.StatusBadRequest)
			return
		}
		// Check bounds before casting
		if recipientID > uint64(math.MaxUint) {
			http.Error(w, "recipient_id out of bounds", http.StatusBadRequest)
			return
		}
		// To fetch 1-on-1 history, we need messages sent FROM user A TO user B,
		// AND messages sent FROM user B TO user A.
		query = query.Where(
			"(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
			userID, uint(recipientID),
			uint(recipientID), userID,
		).Where("team_id IS NULL") // Ensure it's not a group chat message

	} else {
		http.Error(w, "Must provide either team_id or recipient_id", http.StatusBadRequest)
		return
	}

	// 3. Execute the Query
	if err := query.Find(&messages).Error; err != nil {
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}

	// 4. Transform to DTOs and Reverse Order (to get oldest first for display)
	responseDTOs := make([]dto.MessageResponseDto, len(messages))
	for i, msg := range messages {
		// Fill responseDTOs from back to front to get them in ascending (chronological) order
		responseDTOs[len(messages)-1-i] = dto.ToMessageResponseDto(msg)
	}

	// 5. Send Response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseDTOs); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
