package dto

import "time"

// Request DTOs

type EulaAcceptDto struct {
	EulaVersionID uint `json:"eula_version_id" validate:"required"`
}

// Response DTOs

type EulaVersionDto struct {
	ID          uint   `json:"id"`
	Version     string `json:"version"`
	Locale      string `json:"locale"`
	Content     string `json:"content"`
	ContentHash string `json:"content_hash"`
}

type EulaStatusDto struct {
	Locale         string     `json:"locale"`
	ActiveVersion  string     `json:"active_version"`
	Accepted       bool       `json:"accepted"`
	AcceptedAt     *time.Time `json:"accepted_at,omitempty"`
	RequiresAction bool       `json:"requires_action"` // true if user needs to accept
}
