package controllers

import (
	"encoding/json"
	"net/http"
	"server/api/controllers/helpers"
	"server/api/middleware"
	"server/common/appError"
	"server/common/dto"
	"server/common/models"
	"server/common/services"
	"server/common/validator"
)

func CreateEmergencyContact(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	// Decode request
	var req dto.CreateEmergencyInfoDto
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate request
	err = validator.V.Struct(req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to model
	emergencyInfo := dto.CreateEmergencyInfoDtoToModel(req)

	err = services.CreateEmergencyContact(*user, emergencyInfo)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateEmergencyContact(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	emergencyContactId, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Decode request
	var req dto.CreateEmergencyInfoDto
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate request
	err = validator.V.Struct(req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to model
	emergencyInfo := dto.CreateEmergencyInfoDtoToModel(req)

	err = services.UpdateEmergencyContact(*user, emergencyInfo, emergencyContactId)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteEmergencyContact(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	// Get emergency contact id
	emergencyContactId, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.DeleteEmergencyContact(*user, emergencyContactId)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
