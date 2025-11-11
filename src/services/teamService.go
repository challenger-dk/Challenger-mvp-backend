package services

import (
	"server/config"
	"server/models"

	"gorm.io/gorm"
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

	err := config.DB.Preload("Users").Find(&teams).Error
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func GetTeamsByUserId(id uint) ([]models.Team, error) {
	var user models.User

	err := config.DB.Preload("Teams.Users").Preload("Teams.Creator").First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return user.Teams, nil
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

	// Add creator to team users
	t.Users = append(t.Users, creator)

	if err := config.DB.Create(&t).Error; err != nil {
		return models.Team{}, err
	}

	if err := config.DB.Preload("Users").Preload("Creator").First(&t, t.ID).Error; err != nil {
		return models.Team{}, err
	}

	return t, nil
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

// Package private methods
func addUserToTeam(teamId uint, userId uint, db *gorm.DB) error {
	var t models.Team
	var u models.User

	err := db.First(&t, teamId).Error
	if err != nil {
		return err
	}

	err = db.First(&u, userId).Error
	if err != nil {
		return err
	}

	err = db.Model(&t).Association("Users").Append(&u)
	if err != nil {
		return err
	}
	return nil
}
