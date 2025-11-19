package models

import (
	"time"
)

type Sport struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func GetAllowedSports() []string {
	return []string{
		"Tennis",
		"Football",
		"Basketball",
		"PadelTennis",
		"TableTennis",
		"Golf",
		"Volleyball",
		"Badminton",
		"Boxing",
		"Squash",
		"Petanque",
		"Hockey",
		"Handball",
		"Running",
		"Biking",
		"Minigolf",
		"Climbing",
		"Skateboarding",
		"Surfing",
		"Hiking",
		"UltimateFrisbee",
		"Floorball",
	}
}
