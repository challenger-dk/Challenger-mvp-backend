package middleware

import (
	"context"
	"net/http"
	"strings"

	"server/common/appError"
	"server/common/services"
)

// contextKey is an unexported type to avoid key collisions in context.
type contextKey string

// UserContextKey is the key used to store the authenticated user in request context.
const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			appError.HandleError(w, appError.ErrAuthHeaderMissing)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			appError.HandleError(w, appError.ErrInvalidAuthHeader)
			return
		}

		token := parts[1]
		claims, err := services.ValidateJWTToken(token)
		if err != nil {
			appError.HandleError(w, appError.ErrInvalidToken)
			return
		}

		user, err := services.GetUserByID(claims.UserID)
		if err != nil {
			appError.HandleError(w, appError.ErrUserNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
