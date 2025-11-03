package models

import (
    "time"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	Email     string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}