package controllers

import (
	"encoding/json"
	"net/http"
	"server/dto"
	"server/middleware"
	"server/models"
	"server/services"
	"strconv"
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

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Reload user with favorite sports preloaded
	userWithSports, err := services.GetUserByID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToUserResponseDto(*userWithSports))
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	user, err := services.GetUserByID(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto.ToUserResponseDto(*user))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	req := dto.UserUpdateDto{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = services.UpdateUser(uint(id), req)
	if err != nil {
		if err == services.ErrInvalidSport {
			http.Error(w, "One or more invalid sport names provided", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = services.DeleteUser(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
