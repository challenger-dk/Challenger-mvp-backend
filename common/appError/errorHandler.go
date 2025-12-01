package appError

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"server/common/models"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// New struct for structured validation errors
type ValidationErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details"`
}

// --- HandleError (Updated) ---
func HandleError(w http.ResponseWriter, err error) {
	// Log the error using structured logging
	slog.Error("Request Error", "error", err)

	// 1. Check for validation errors
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		w.WriteHeader(http.StatusBadRequest)

		// Build the structured error details
		details := make(map[string]string)
		for _, fe := range validationErrors {
			// fe.Namespace() returns the full path, e.g., "ChallengeCreateDto.location.address"
			fieldName := fe.Namespace()

			// We want to remove the struct name ("ChallengeCreateDto.") to get just "location.address"
			if dotIndex := strings.Index(fieldName, "."); dotIndex != -1 {
				fieldName = fieldName[dotIndex+1:]
			}

			// If stripping failed or there was no dot, fallback to the simple field name
			if fieldName == "" {
				fieldName = fe.Field()
			}

			details[fieldName] = getValidationErrorMessage(fe)
		}

		resp := ValidationErrorResponse{
			Error:   "Validation Failed",
			Details: details,
		}

		if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
			slog.Error("Failed to encode validation error response", "error", encodeErr)
		}
		return
	}

	// 2. Iterate the error map
	for statusCode, knownErrors := range errorMap {
		for _, knownErr := range knownErrors {
			if errors.Is(err, knownErr) {
				w.WriteHeader(statusCode)
				var resp ErrorResponse
				if knownErr == gorm.ErrRecordNotFound {
					resp = ErrorResponse{Error: "Resource not found"}
				} else {
					resp = ErrorResponse{Error: err.Error()}
				}
				if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
					slog.Error("Failed to encode error response", "error", encodeErr)
				}
				return
			}
		}
	}

	// 3. Default to 500
	w.WriteHeader(http.StatusInternalServerError)
	if encodeErr := json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("An unexpected error occured: %s", err.Error())}); encodeErr != nil {
		slog.Error("Failed to encode default error response", "error", encodeErr)
	}
}

// --- New helper function to create clean error messages ---
func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "This must be a valid email address"
	case "min":
		return fmt.Sprintf("This field must be at least %s characters long", fe.Param())
	case "max":
		return fmt.Sprintf("This field must be no more than %s characters long", fe.Param())
	case "latitude":
		return "This must be a valid latitude (between -90 and 90)"
	case "longitude":
		return "This must be a valid longitude (between -180 and 180)"
	case "is-valid-sport":
		return fmt.Sprintf("Not valid sport should be one of: [%s]", strings.Join(models.GetAllowedSports(), ", "))
	case "oneof":
		return fmt.Sprintf("This field must be one of: %s", fe.Param())
	default:
		return fmt.Sprintf("Field validation error message not supported for: %s", fe.Tag())
	}
}
