package middleware

import (
	"context"
	"net/http"
	"strings"

	"server/services"
)

// contextKey is an unexported type to avoid key collisions in context.
type contextKey string

// UserContextKey is the key used to store the authenticated user in request context.
const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		claims, err := services.ValidateJWTToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		user, err := services.GetUserByID(claims.UserID)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func OptionalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]
				claims, err := services.ValidateJWTToken(token)
				if err == nil {
					user, err := services.GetUserByID(claims.UserID)
					if err == nil {
						ctx := context.WithValue(r.Context(), UserContextKey, user)
						r = r.WithContext(ctx)
					}
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
