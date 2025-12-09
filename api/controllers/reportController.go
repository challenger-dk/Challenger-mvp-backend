package controllers

import (
	"encoding/json"
	"net/http"
	"server/api/middleware"
	"server/common/appError"
	"server/common/dto"
	"server/common/models"
	"server/common/services"
	"server/common/validator"
)

func CreateReport(w http.ResponseWriter, r *http.Request) {
	// 1. Get current user
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	// 2. Decode request
	var req dto.ReportCreateDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appError.HandleError(w, err)
		return
	}

	// 3. Validate
	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	// 4. Save to DB
	if err := services.CreateReport(user.ID, req); err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
