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

func GetTeams() ([]models.Team, error) {
	var teams []models.Team

	err := config.DB.Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func CreateTeam(team dto.TeamCreateDto) (dto.TeamResponseDto, error) {
	t := dto.TeamCreateDtoToModel(team)

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

func DeleteTeam(id uint) error {
	var t models.Team

	err := config.DB.First(&t, id).Error
	if err != nil {
		return err
	}

	// Remove user associations
	err = config.DB.Model(&t).Association("Users").Clear()
	if err != nil {
		return err
	}

	err = config.DB.Delete(&t).Error
	if err != nil {
		return err
	}

	return nil
}
