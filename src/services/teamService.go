package services

import (
	"server/config"
	"server/dto"
	"server/models"
)

func GetTeamByID(id uint) (dto.TeamResponseDto, error) {
	var t models.Team

	err := config.DB.Preload("Users").First(&t, id).Error
	if err != nil {
		return dto.TeamResponseDto{}, err
	}

	return dto.ToTeamResponseDto(t), nil
}

func GetTeams() ([]models.User, error) {
	var users []models.User

	err := config.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func CreateTeam(user dto.TeamCreateDto) (dto.TeamResponseDto, error) {
	t := dto.TeamCreateDtoToModel(user)

	err := config.DB.Create(&t).Error
	if err != nil {
		return dto.TeamResponseDto{}, err
	}

	resp := dto.ToTeamResponseDto(t)

	return resp, nil
}

func UpdateTeam(id uint, team dto.TeamCreateDto) error {
	var t models.Team

	err := config.DB.First(&t, id).Error
	if err != nil {
		return err
	}

	if team.Name != "" {
		t.Name = team.Name
	}

	err = config.DB.Save(&t).Error
	if err != nil {
		return err
	}

	return nil
}
