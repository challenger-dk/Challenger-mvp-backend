package models

import "time"

type ReportTargetType string

const (
	ReportTargetUser      ReportTargetType = "USER"
	ReportTargetTeam      ReportTargetType = "TEAM"
	ReportTargetChallenge ReportTargetType = "CHALLENGE"
	ReportTargetMessage   ReportTargetType = "MESSAGE"
)

type Report struct {
	ID         uint `gorm:"primaryKey"`
	ReporterID uint `gorm:"not null"` // Who sent the report
	Reporter   User `gorm:"foreignKey:ReporterID"`

	TargetID   uint             `gorm:"not null"` // ID of the thing being reported
	TargetType ReportTargetType `gorm:"not null"` // "USER", "TEAM", etc.

	Reason  string `gorm:"not null"` // "RACISM", "SPAM", etc.
	Comment string

	Status    string `gorm:"default:'PENDING'"` // PENDING, RESOLVED, DISMISSED
	CreatedAt time.Time
}
