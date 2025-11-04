package controllers

import (
	"encoding/json"
	"net/http"
	"server/dto"
	"server/services"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := services.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	response := make([]dto.UserResponseDto, len(users))
	for i, user := range users {
		response[i] = dto.ToUserResponseDto(user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
