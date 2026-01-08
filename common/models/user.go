package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID             uint    `gorm:"primaryKey"`
	Email          string  `gorm:"not null;unique"`
	Password       *string `gorm:"default:null"` // Nullable for OAuth users
	AuthProvider   string  `gorm:"default:''"`   // "google", "apple", or empty for regular users
	FirstName      string  `gorm:"not null"`
	LastName       string
	ProfilePicture string
	Bio            string
	BirthDate      time.Time
	City           string

	// Password Reset
	PasswordResetCode          string     `gorm:"index"`
	PasswordResetCodeExpiresAt *time.Time `gorm:"index"`

	// Relationships
	FavoriteSports    []Sport     `gorm:"many2many:user_favorite_sports;"`
	Teams             []Team      `gorm:"many2many:user_teams;"`
	CreatedChallenges []Challenge `gorm:"foreignKey:CreatorID"`
	JoinedChallenges  []Challenge `gorm:"many2many:user_challenges;"`
	Friends           []User      `gorm:"many2many:user_friends;joinForeignKey:UserID;JoinReferences:FriendID"`
	BlockedUsers      []User      `gorm:"many2many:user_blocked_users;joinForeignKey:UserID;JoinReferences:BlockedUserID"`

	Settings          *UserSettings   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	EmergencyContacts []EmergencyInfo `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	//Other
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
