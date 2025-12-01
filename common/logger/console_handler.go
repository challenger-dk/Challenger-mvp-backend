package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

// ConsoleHandler is a slog.Handler that writes colored, human-readable logs
type ConsoleHandler struct {
	w  io.Writer
	mu sync.Mutex
}

func NewConsoleHandler(w io.Writer) *ConsoleHandler {
	return &ConsoleHandler{w: w}
}

func (h *ConsoleHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h // Simplified for brevity; normally you'd clone and append
}

func (h *ConsoleHandler) WithGroup(name string) slog.Handler {
	return h // Simplified
}

func (h *ConsoleHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 1. Colorize Level
	level := r.Level.String()
	color := "\033[0m" // Reset

	switch r.Level {
	case slog.LevelDebug:
		level = "DBG"
		color = "\033[36m" // Cyan
	case slog.LevelInfo:
		level = "INF"
		color = "\033[32m" // Green
	case slog.LevelWarn:
		level = "WRN"
		color = "\033[33m" // Yellow
	case slog.LevelError:
		level = "ERR"
		color = "\033[31m" // Red
	}

	// 2. Format Timestamp
	timeStr := r.Time.Format(time.TimeOnly)

	// 3. Write Header: Time | Level | Message
	// Example: 12:48:51 | WRN | HTTP Request
	fmt.Fprintf(h.w, "%s \033[90m|\033[0m %s%s\033[0m \033[90m|\033[0m %-15s",
		timeStr, color, level, r.Message)

	// 4. Format Attributes
	r.Attrs(func(a slog.Attr) bool {
		// Skip empty fields
		if a.Value.Kind() == slog.KindString && a.Value.String() == "" {
			return true
		}

		valStr := a.Value.String()

		// Special formatting for specific keys
		switch a.Key {
		case "status":
			// Colorize HTTP status codes
			status := int(a.Value.Int64())
			if status >= 500 {
				valStr = fmt.Sprintf("\033[31m%d\033[0m", status) // Red
			} else if status >= 400 {
				valStr = fmt.Sprintf("\033[33m%d\033[0m", status) // Yellow
			} else {
				valStr = fmt.Sprintf("\033[32m%d\033[0m", status) // Green
			}
		case "error":
			valStr = fmt.Sprintf("\033[31m%s\033[0m", valStr)
		case "duration":

		}

		fmt.Fprintf(h.w, " \033[90m%s=\033[0m%s", a.Key, valStr)
		return true
	})

	fmt.Fprintln(h.w)
	return nil
}
