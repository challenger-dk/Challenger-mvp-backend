package middleware

import (
	"log"
	"net/http"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[92m" // Bright Green
	ColorYellow = "\033[93m" // Bright Yellow
	ColorBlue   = "\033[94m" // Bright Blue
	ColorRed    = "\033[91m" // Bright Red
	ColorCyan   = "\033[96m" // Bright Cyan
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		var methodColor string
		switch r.Method {
		case http.MethodGet:
			methodColor = ColorGreen
		case http.MethodPost:
			methodColor = ColorYellow
		case http.MethodPut:
			methodColor = ColorBlue
		case http.MethodDelete:
			methodColor = ColorRed
		default:
			methodColor = ColorCyan
		}

		log.Printf(
			"%s[%s]%s %s (%.6f s)",
			methodColor,
			r.Method,
			ColorReset,
			r.RequestURI,
			time.Since(start).Seconds(),
		)
	})
}
