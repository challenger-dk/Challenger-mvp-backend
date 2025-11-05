package services

import (
	"server/config"
	"server/models"
)

func GetChallengeByID(id uint) (models.Challenge, error) {
	var c models.Challenge

	err := config.DB.Preload("Users").Preload("Teams").Preload("Creator").First(&c, id).Error
	if err != nil {
		return models.Challenge{}, err
	}

	return c, nil
}

// TODO: Some kind of pagination so we dont fetch all challenges
func GetChallenges() ([]models.Challenge, error) {
	var challenges []models.Challenge

	err := config.DB.Preload("Users").Preload("Teams").Preload("Creator").Find(&challenges).Error
	if err != nil {
		return nil, err
	}

	return challenges, nil
}

func CreateChallenge(c models.Challenge) (models.Challenge, error) {
	// ensure creator exists
	creator := models.User{}
	if err := config.DB.First(&creator, c.CreatorID).Error; err != nil {
		return models.Challenge{}, err
	}

	c.CreatorID = creator.ID
	c.Creator = models.User{}

	if err := config.DB.Create(&c).Error; err != nil {
		return models.Challenge{}, err
	}

	if err := config.DB.Preload("Users").Preload("Teams").Preload("Creator").First(&c, c.ID).Error; err != nil {
		return models.Challenge{}, err
	}

	return c, nil
}

func UpdateChallenge(id uint, ch models.Challenge) error {
	var c models.Challenge

	err := config.DB.First(&c, id).Error
	if err != nil {
		return err
	}

	if ch.Name != "" {
		c.Name = ch.Name
	}

	if ch.Description != "" {
		c.Description = ch.Description
	}

	if ch.Sport != "" {
		c.Sport = ch.Sport
	}

	if ch.Location != "" {
		c.Location = ch.Location
	}

	err = config.DB.Save(&c).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteChallenge(id uint) error {
	var c models.Challenge

	err := config.DB.First(&c, id).Error
	if err != nil {
		return err
	}

	// clear many2many associations to avoid orphan join rows
	err = config.DB.Model(&c).Association("Teams").Clear()
	if err != nil {
		return err
	}

	err = config.DB.Model(&c).Association("Users").Clear()
	if err != nil {
		return err
	}

	err = config.DB.Delete(&c).Error
	if err != nil {
		return err
	}
	return nil
}
