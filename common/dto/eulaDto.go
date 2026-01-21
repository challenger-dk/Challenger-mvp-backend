package dto

import "time"

// Request DTOs

type EulaAcceptDto struct {
	EulaVersionID uint `json:"eula_version_id" validate:"required"`
}

// Response DTOs

type EulaVersionDto struct {
	ID          uint   `json:"id" validate:"sanitize"`
	Version     string `json:"version" validate:"sanitize"`
	Locale      string `json:"locale" validate:"sanitize"`
	Content     string `json:"content" validate:"sanitize"`
	ContentHash string `json:"content_hash" validate:"sanitize"`
}

type EulaStatusDto struct {
	Locale         string     `json:"locale" validate:"sanitize"`
	ActiveVersion  string     `json:"active_version" validate:"sanitize"`
	Accepted       bool       `json:"accepted"`
	AcceptedAt     *time.Time `json:"accepted_at,omitempty"`
	RequiresAction bool       `json:"requires_action"` // true if user needs to accept
}
