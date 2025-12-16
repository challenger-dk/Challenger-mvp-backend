package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func SlogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		status := ww.Status()

		// 1. Determine Log Level based on Status Code
		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		// 2. Prepare Attributes
		attrs := []any{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", status),
			slog.String("status_text", http.StatusText(status)), // <--- Adds "Method Not Allowed"
			slog.Duration("duration", duration),
		}

		// Only add these if they add value (reduce noise)
		reqID := middleware.GetReqID(r.Context())
		if reqID != "" {
			attrs = append(attrs, slog.String("req_id", reqID))
		}

		// 3. Log with dynamic level
		slog.Log(r.Context(), level, "HTTP Request", attrs...)
	})
}
