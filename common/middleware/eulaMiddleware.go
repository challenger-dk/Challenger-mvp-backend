package middleware

import (
	"net/http"
	"server/common/appError"
	"server/common/models"
	"server/common/services"
)

// EulaMiddleware checks if the authenticated user has accepted the active EULA
// This middleware should be applied AFTER AuthMiddleware
func EulaMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user from context (set by AuthMiddleware)
		user, ok := r.Context().Value(UserContextKey).(*models.User)
		if !ok {
			// If no user in context, let it pass (AuthMiddleware will handle it)
			next.ServeHTTP(w, r)
			return
		}

		// Get locale from query parameter or header (default: da-DK)
		locale := r.URL.Query().Get("locale")
		if locale == "" {
			locale = r.Header.Get("Accept-Language")
		}
		if locale == "" {
			locale = "da-DK"
		}

		// Check if user has accepted the active EULA
		accepted, err := services.HasUserAcceptedActiveEula(user.ID, locale)
		if err != nil {
			// If there's an error checking EULA status, log it but don't block
			// (This prevents EULA system from breaking the entire app)
			// In production, you might want to handle this differently
			next.ServeHTTP(w, r)
			return
		}

		if !accepted {
			// User has not accepted the active EULA
			appError.HandleError(w, appError.ErrEulaNotAccepted)
			return
		}

		// User has accepted EULA, proceed
		next.ServeHTTP(w, r)
	})
}
