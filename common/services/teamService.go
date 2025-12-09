package services

import (
	"server/common/appError"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

// --- GET ---
func GetTeamByID(id uint) (models.Team, error) {
	var t models.Team

	err := config.DB.Preload("Users").
		Preload("Creator").
		Preload("Location").
		First(&t, id).
		Error

	if err != nil {
		return models.Team{}, err
	}

	return t, nil
}

func GetTeams() ([]models.Team, error) {
	var teams []models.Team

	err := config.DB.Preload("Users").
		Preload("Creator").
		Preload("Location").
		Find(&teams).
		Error

	if err != nil {
		return nil, err
	}
	return teams, nil
}

func GetTeamsByUserId(id uint) ([]models.Team, error) {
	var user models.User

	err := config.DB.Preload("Teams.Users").
		Preload("Teams.Creator").
		Preload("Teams.Location").
		First(&user, id).
		Error

	if err != nil {
		return nil, err
	}

	return user.Teams, nil
}

// --- POST ---
func CreateTeam(t models.Team) (models.Team, error) {
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if t.Location != nil {
			location, err := FindOrCreateLocation(tx, *t.Location)
			if err != nil {
				return err
			}

			t.LocationID = &location.ID
		}

		// Set Location to nil to avoid duplicate create
		t.Location = nil

		creator := models.User{}
		err := tx.First(&creator, t.CreatorID).Error
		if err != nil {
			return err
		}

		t.CreatorID = creator.ID
		t.Creator = models.User{}
		t.Users = append(t.Users, creator)

		err = tx.Create(&t).Error
		if err != nil {
			return err
		}

		err = tx.Preload("Users").
			Preload("Creator").
			Preload("Location").
			First(&t, t.ID).
			Error

		return err
	})

	if err != nil {
		return models.Team{}, err
	}

	return t, nil
}

// --- PUT ---
func UpdateTeam(id uint, team models.Team) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var t models.Team

		err := tx.First(&t, id).Error
		if err != nil {
			return err
		}

		if team.Name != "" {
			t.Name = team.Name
		}

		return tx.Save(&t).Error
	})
}

// --- DELETE ---
func DeleteTeam(id uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var t models.Team

		// Load team + users in a single query
		if err := tx.Preload("Users").First(&t, id).Error; err != nil {
			return err
		}

		// Users to notify
		users := t.Users
		creatorId := t.CreatorID

		// Soft delete the team.
		// With gorm.DeletedAt on Team, this sets deleted_at instead of hard-deleting.
		if err := tx.Delete(&t).Error; err != nil {
			return err
		}

		// Notify users (except creator)
		for _, u := range users {
			if u.ID == creatorId {
				continue
			}

			CreateTeamDeletedNotification(tx, u, t)
		}

		return nil
	})
}

func RemoveUserFromTeam(creator models.User, teamId uint, userId uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var t models.Team
		var u models.User

		err := tx.First(&t, teamId).Error
		if err != nil {
			return err
		}

		if creator.ID != t.CreatorID {
			return appError.ErrUnauthorized
		}

		err = tx.First(&u, userId).Error
		if err != nil {
			return err
		}

		err = tx.Model(&t).
			Association("Users").
			Delete(&u)

		if err != nil {
			return err
		}

		// Notification
		CreateRemovedUserFromTeamNotification(tx, u.ID, t)

		return nil
	})
}

func LeaveTeam(user models.User, teamId uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var t models.Team

		err := tx.First(&t, teamId).Error
		if err != nil {
			return err
		}

		err = tx.Model(&t).
			Association("Users").
			Delete(&user)

		if err != nil {
			return err
		}

		//Notification
		CreateUserLeftTeamNotification(tx, user, t)

		return nil
	})
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

	err = db.Model(&t).
		Association("Users").
		Append(&u)

	if err != nil {
		return err
	}

	return nil
}
