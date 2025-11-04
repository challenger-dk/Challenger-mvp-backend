package services

import (
	"server/config"
	"server/models"
)

func GetAllSports() ([]models.Sport, error) {
	var sports []models.Sport
	if err := config.DB.Find(&sports).Error; err != nil {
		return nil, err
	}
	return sports, nil
}

func SeedSports() error {
	allowedSports := models.GetAllowedSports()

	for _, sportName := range allowedSports {
		var sport models.Sport
		if err := config.DB.Where("name = ?", sportName).FirstOrCreate(&sport, models.Sport{Name: sportName}).Error; err != nil {
			return err
		}
	}

	return nil
}
