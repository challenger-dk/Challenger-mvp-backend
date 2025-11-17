package controllers

import (
	"encoding/json"
	"net/http"
	"server/appError"
	"server/controllers/helpers"
	"server/dto"
	"server/middleware"
	"server/models"
	"server/services"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := services.GetUsers()
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to response DTOs
	response := make([]dto.UserResponseDto, len(users))
	for i, user := range users {
		response[i] = dto.ToUserResponseDto(user)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Reload user with favorite sports preloaded
	userWithSports, err := services.GetUserByID(user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(dto.ToUserResponseDto(*userWithSports))
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	user, err := services.GetUserByID(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(dto.ToUserResponseDto(*user))
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	req := dto.UserUpdateDto{}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.UpdateUser(id, req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.DeleteUser(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
