package controllers

import (
	"encoding/json"
	"net/http"

	"server/common/appError"
	"server/common/dto"
	"server/common/services"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserCreateDto

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	user, err := services.CreateUser(dto.UserCreateDtoToModel(req), req.Password)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	token, err := services.GenerateJWTToken(user)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]any{
		"user":  dto.ToUserResponseDto(*user),
		"token": token,
	})

	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req dto.Login

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	user, token, err := services.Login(req.Email, req.Password)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]any{
		"user":  dto.ToUserResponseDto(*user),
		"token": token,
	})

	if err != nil {
		appError.HandleError(w, err)
		return
	}
}
