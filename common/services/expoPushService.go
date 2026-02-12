package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"server/common/appError"
)

const expoPushAPIURL = "https://exp.host/--/api/v2/push/send"

// ExpoPushMessage represents a single push notification message for the Expo API.
type ExpoPushMessage struct {
	To    string         `json:"to"`
	Title string         `json:"title,omitempty"`
	Body  string         `json:"body,omitempty"`
	Data  map[string]any `json:"data,omitempty"`
}

// ExpoPushTicket represents the response for a single pushed message.
type ExpoPushTicket struct {
	Status  string         `json:"status"`
	ID      string         `json:"id,omitempty"`
	Message string         `json:"message,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

// ExpoPushResponse is the full response from the Expo push API.
type ExpoPushResponse struct {
	Data   []ExpoPushTicket `json:"data,omitempty"`
	Errors []ExpoAPIError   `json:"errors,omitempty"`
}

// ExpoAPIError represents an error from the Expo API (request-level error).
type ExpoAPIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrExpoPush returned when the Expo API reports an error (e.g. invalid token, rate limit).
type ErrExpoPush struct {
	Message string
	Details map[string]any
}

func (e ErrExpoPush) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("expo push error: %s (details: %v)", e.Message, e.Details)
	}
	return fmt.Sprintf("expo push error: %s", e.Message)
}

// isValidExpoPushToken validates that the token follows the ExponentPushToken or ExpoPushToken format.
func isValidExpoPushToken(token string) bool {
	token = strings.TrimSpace(token)
	if token == "" {
		return false
	}
	return strings.HasPrefix(token, "ExponentPushToken[") ||
		strings.HasPrefix(token, "ExpoPushToken[")
}

// SendExpoPushNotification sends a push notification via the Expo Push API.
// It accepts the push token, notification title, body, and optional custom data payload.
// Returns an error if the token is invalid, the HTTP request fails, or the Expo API reports an error.
func SendExpoPushNotification(pushToken string, title string, body string, data map[string]any) error {
	if !isValidExpoPushToken(pushToken) {
		return appError.ErrInvalidPushToken
	}

	message := ExpoPushMessage{
		To:    pushToken,
		Title: title,
		Body:  body,
		Data:  data,
	}

	bodyBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal push message: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, expoPushAPIURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Accept-Encoding", "gzip, deflate")

	client := &http.Client{Timeout: 15 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send push notification: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	return parseExpoPushResponse(response.StatusCode, responseBody)
}

// parseExpoPushResponse interprets the Expo API response and returns an error if delivery failed.
func parseExpoPushResponse(statusCode int, body []byte) error {
	var apiResponse ExpoPushResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return fmt.Errorf("failed to parse expo response (status %d): %w", statusCode, err)
	}

	// Request-level errors (entire request failed)
	if len(apiResponse.Errors) > 0 {
		first := apiResponse.Errors[0]
		return ErrExpoPush{
			Message: first.Message,
			Details: map[string]any{"code": first.Code},
		}
	}

	// Check for 4xx/5xx status codes
	if statusCode >= 400 {
		return fmt.Errorf("expo push API returned status %d: %s", statusCode, string(body))
	}

	// Message-level errors (individual ticket failed)
	if len(apiResponse.Data) > 0 {
		ticket := apiResponse.Data[0]
		if ticket.Status == "error" {
			details := ticket.Details
			if details == nil {
				details = make(map[string]any)
			}
			if ticket.Message != "" {
				details["message"] = ticket.Message
			}
			if errCode, ok := details["error"].(string); ok && errCode == "DeviceNotRegistered" {
				return ErrExpoPush{
					Message: "device is not registered for push notifications",
					Details: details,
				}
			}
			return ErrExpoPush{
				Message: ticket.Message,
				Details: details,
			}
		}
	}

	return nil
}

// IsDeviceNotRegistered returns true if the error indicates the push token is no longer valid.
// Callers should stop sending notifications to this token and optionally clear it from the user record.
func IsDeviceNotRegistered(err error) bool {
	var expoErr ErrExpoPush
	if errors.As(err, &expoErr) {
		if errCode, ok := expoErr.Details["error"].(string); ok {
			return errCode == "DeviceNotRegistered"
		}
	}
	return false
}
