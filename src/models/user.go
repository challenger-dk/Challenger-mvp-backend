package models

import (
	"time"
)

type User struct {
	ID             uint   `gorm:"primaryKey"`
	Email          string `gorm:"not null;unique"`
	Password       string `gorm:"not null"`
	FirstName      string `gorm:"not null"`
	LastName       string
	ProfilePicture string
	Bio            string

	// Relationships
	FavoriteSports    []Sport     `gorm:"many2many:user_favorite_sports;"`
	Teams             []Team      `gorm:"many2many:user_teams;"`
	CreatedChallenges []Challenge `gorm:"foreignKey:CreatorID"`
	JoinedChallenges  []Challenge `gorm:"many2many:user_challenges;"`
	Friends           []User      `gorm:"many2many:user_friends;joinForeignKey:UserID;JoinReferences:FriendID"`
	CreatedAt         time.Time   `gorm:"autoCreateTime"`
	UpdatedAt         time.Time   `gorm:"autoUpdateTime"`
}
