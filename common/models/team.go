package models

import (
	"time"

	"gorm.io/gorm"
)

type TeamRole string

const (
	RoleMember TeamRole = "member"
	RoleAdmin  TeamRole = "admin"
	RoleOwner  TeamRole = "owner"
)

type Team struct {
	ID uint `gorm:"primaryKey"`

	Name        string `gorm:"not null"`
	Description *string

	Sports     []Sport      `gorm:"many2many:team_sports;"`
	Users      []TeamMember `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Creator    User         `gorm:"foreignKey:CreatorID"`
	CreatorID  uint         `gorm:"not null"`
	LocationID *uint        `gorm:"index"`
	Location   *Location    `gorm:"foreignKey:LocationID"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type TeamMember struct {
	TeamID uint `gorm:"primaryKey;index"`
	UserID uint `gorm:"primaryKey;index"`

	Role TeamRole `gorm:"type:varchar(20);not null;default:'member';index"`

	Team      Team      `gorm:"foreignKey:TeamID"`
	User      User      `gorm:"foreignKey:UserID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
