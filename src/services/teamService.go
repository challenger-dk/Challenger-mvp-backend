package services

import (
	"server/config"
	"server/models"
)

// --- GET ---
func GetTeamByID(id uint) (models.Team, error) {
	var t models.Team

	err := config.DB.Preload("Users").Preload("Creator").First(&t, id).Error
	if err != nil {
		return models.Team{}, err
	}

	return t, nil
}

func GetTeams() ([]models.Team, error) {
	var teams []models.Team

	err := config.DB.Preload("Users").Preload("Creator").Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func GetTeamsByUserId(id uint) ([]models.Team, error) {
	var teams []models.Team

	// Find teams where CreatorID == id
	err := config.DB.Preload("Users").Preload("Creator").
		Where("creator_id = ?", id).
		Find(&teams).Error

	if err != nil {
		return nil, err
	}

	return teams, nil
}

// --- POST ---
func CreateTeam(t models.Team) (models.Team, error) {
	// Ensure creator exists
	creator := models.User{}
	if err := config.DB.First(&creator, t.CreatorID).Error; err != nil {
		return models.Team{}, err
	}

	// Set foreign key explicitly and avoid nested create/update
	t.CreatorID = creator.ID
	t.Creator = models.User{}

	if err := config.DB.Create(&t).Error; err != nil {
		return models.Team{}, err
	}

	if err := config.DB.Preload("Users").Preload("Creator").First(&t, t.ID).Error; err != nil {
		return models.Team{}, err
	}

	return t, nil
}

func AddUserToTeam(teamId uint, userId uint) error {
	var t models.Team
	var u models.User

	err := config.DB.First(&t, teamId).Error
	if err != nil {
		return err
	}

	err = config.DB.First(&u, userId).Error
	if err != nil {
		return err
	}

	err = config.DB.Model(&t).Association("Users").Append(&u)
	if err != nil {
		return err
	}
	return nil
}

// --- PUT ---
func UpdateTeam(id uint, team models.Team) error {
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

// --- DELETE ---
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
