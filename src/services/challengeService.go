package services

import (
	"server/config"
	"server/dto"
	"server/models"
)

func GetChallengeByID(id uint) (dto.ChallengeResponseDto, error) {
	var c models.Challenge

	err := config.DB.Preload("Users").Preload("Teams").Preload("Creator").First(&c, id).Error
	if err != nil {
		return dto.ChallengeResponseDto{}, err
	}

	return dto.ToChallengeResponseDto(c), nil
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

func CreateChallenge(ch dto.ChallengeCreateDto) (dto.ChallengeResponseDto, error) {
	c := dto.ChallengeCreateDtoToModel(ch)

	creator, err := GetUserByID(ch.CreatorId)
	if err != nil {
		return dto.ChallengeResponseDto{}, err
	}

	c.Creator = *creator

	err = config.DB.Create(&c).Error
	if err != nil {
		return dto.ChallengeResponseDto{}, err
	}

	return dto.ToChallengeResponseDto(c), nil
}

func UpdateChallenge(id uint, ch dto.ChallengeCreateDto) error {
	var c models.Challenge

	if err := config.DB.First(&c, id).Error; err != nil {
		return err
	}

	if ch.Name != "" {
		c.Name = ch.Name
	}
	// extend here if ChallengeCreateDto grows (description, sport, etc.)

	if err := config.DB.Save(&c).Error; err != nil {
		return err
	}
	return nil
}

func DeleteChallenge(id uint) error {
	var c models.Challenge

	if err := config.DB.First(&c, id).Error; err != nil {
		return err
	}

	// clear many2many associations to avoid orphan join rows
	if err := config.DB.Model(&c).Association("Teams").Clear(); err != nil {
		return err
	}
	if err := config.DB.Model(&c).Association("Users").Clear(); err != nil {
		return err
	}

	if err := config.DB.Delete(&c).Error; err != nil {
		return err
	}
	return nil
}
