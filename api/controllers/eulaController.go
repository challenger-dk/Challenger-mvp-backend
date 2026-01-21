package controllers

import (
	"encoding/json"
	"net/http"
	"server/common/appError"
	"server/common/dto"
	"server/common/middleware"
	"server/common/models"
	"server/common/services"
	"server/common/validator"
)

// GetActiveEula returns the active EULA for the specified locale
func GetActiveEula(w http.ResponseWriter, r *http.Request) {
	// Get locale from query parameter (default: da-DK)
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "da-DK"
	}

	eulaVersion, err := services.GetActiveEula(locale)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to DTO
	response := dto.EulaVersionDto{
		ID:          eulaVersion.ID,
		Version:     eulaVersion.Version,
		Locale:      eulaVersion.Locale,
		Content:     eulaVersion.Content,
		ContentHash: eulaVersion.ContentHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetEulaStatus checks if the current user has accepted the active EULA
func GetEulaStatus(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	// Get locale from query parameter (default: da-DK)
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "da-DK"
	}

	// Get active EULA for locale
	activeEula, err := services.GetActiveEula(locale)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Check if user has accepted this version
	acceptance, err := services.GetUserEulaAcceptance(user.ID, activeEula.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Build response DTO
	status := dto.EulaStatusDto{
		Locale:         locale,
		ActiveVersion:  activeEula.Version,
		Accepted:       acceptance != nil,
		AcceptedAt:     nil,
		RequiresAction: acceptance == nil,
	}

	if acceptance != nil {
		status.AcceptedAt = &acceptance.AcceptedAt
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// AcceptEula records the user's acceptance of a EULA version
func AcceptEula(w http.ResponseWriter, r *http.Request) {
	// Get current user
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	// Decode request
	var req dto.EulaAcceptDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	// Accept EULA
	if err := services.AcceptEula(user.ID, req.EulaVersionID); err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
