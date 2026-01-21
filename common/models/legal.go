package models

import (
	"time"
)

type EulaVersion struct {
	ID          uint      `gorm:"primaryKey"`
	Version     string    `gorm:"not null;uniqueIndex:idx_eula_version_locale"`
	Locale      string    `gorm:"not null;uniqueIndex:idx_eula_version_locale"`
	Content     string    `gorm:"type:text;not null"`
	ContentHash string    `gorm:"type:char(64);not null"`
	IsActive    bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type EulaAcceptance struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        uint      `gorm:"not null;uniqueIndex:idx_user_eula"`
	EulaVersionID uint      `gorm:"not null;uniqueIndex:idx_user_eula"`
	AcceptedAt    time.Time `gorm:"not null;autoCreateTime"`
}
