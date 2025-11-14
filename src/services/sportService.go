package services

import (
	"server/config"
	"server/models"
)

func GetAllSports() ([]models.Sport, error) {
	var sports []models.Sport

	err := config.DB.Find(&sports).Error
	if err != nil {
		return nil, err
	}

	return sports, nil
}

// Package-level
// associateFavoriteSports validates and associates sports with a user
func associateFavoriteSports(userID uint, sportNames []string) error {
	// Validate sport names against the global cache
	for _, sportName := range sportNames {
		if _, ok := config.SportsCache[sportName]; !ok {
			return ErrInvalidSport
		}
	}

	// Find or create sports in database
	var sports []models.Sport
	for _, sportName := range sportNames {
		var sport models.Sport

		err := config.DB.Where("name = ?", sportName).
			FirstOrCreate(&sport, models.Sport{Name: sportName}).
			Error

		if err != nil {
			return err
		}

		sports = append(sports, sport)
	}

	// Replace user's favorite sports
	var user models.User

	err := config.DB.First(&user, userID).Error
	if err != nil {
		return err
	}

	return config.DB.Model(&user).
		Association("FavoriteSports").
		Replace(sports)
}
