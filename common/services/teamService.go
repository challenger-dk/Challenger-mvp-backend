package services

import (
	"server/common/appError"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

// --- GET ---
func GetTeamByID(id uint, currentUserID uint) (models.Team, error) {
	var t models.Team

	err := config.DB.
		Scopes(ExcludeBlockedUsersOn(currentUserID, "creator_id")).
		Preload("Users.User", ExcludeBlockedUsers(currentUserID)).
		Preload("Creator").
		Preload("Location").
		First(&t, id).
		Error

	if err != nil {
		return models.Team{}, err
	}

	return t, nil
}

func GetTeams(currentUserID uint) ([]models.Team, error) {
	var teams []models.Team

	err := config.DB.
		Scopes(ExcludeBlockedUsersOn(currentUserID, "creator_id")).
		Preload("Users.User", ExcludeBlockedUsers(currentUserID)).
		Preload("Creator").
		Preload("Location").
		Find(&teams).
		Error

	if err != nil {
		return nil, err
	}
	return teams, nil
}

func GetTeamsByUserId(id uint, currentUserID uint) ([]models.TeamMember, error) {
	var user models.User

	err := config.DB.
		Preload("Teams.Team", ExcludeBlockedUsersOn(currentUserID, "creator_id")).
		Preload("Teams.Team.Users.User", ExcludeBlockedUsers(currentUserID)).
		Preload("Teams.Team.Creator").
		Preload("Teams.Team.Location").
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
		t.Users = append(t.Users, models.TeamMember{
			UserID: creator.ID,
			Role:   models.RoleOwner,
		})

		err = tx.Create(&t).Error
		if err != nil {
			return err
		}

		err = tx.Preload("Users.User").
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
// Soft delete team and associations (soft delete)
func SoftDeleteTeam(id uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var t models.Team

		// Load team + users in a single query
		if err := tx.Preload("Users.User").First(&t, id).Error; err != nil {
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
			if u.UserID == creatorId {
				continue
			}

			CreateTeamDeletedNotification(tx, u.User, t)
		}

		return nil
	})
}

// Completely Delete team and associations (no soft delete)
func DeleteTeam(id uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var t models.Team
		// Load team + users in a single query
		if err := tx.Preload("Users.User").First(&t, id).Error; err != nil {
			return err
		}

		// Delete team
		if err := tx.Unscoped().Delete(&t).Error; err != nil {
			return err
		}

		return nil
	})
}

func RemoveUserFromTeam(creator models.User, teamId uint, userId uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var t models.Team

		err := tx.First(&t, teamId).Error
		if err != nil {
			return err
		}

		if creator.ID != t.CreatorID {
			return appError.ErrUnauthorized
		}

		// Delete the TeamMember record
		err = tx.Where("team_id = ? AND user_id = ?", teamId, userId).
			Delete(&models.TeamMember{}).Error

		if err != nil {
			return err
		}

		// Notification
		CreateRemovedUserFromTeamNotification(tx, userId, t)

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

		// Delete the TeamMember record
		err = tx.Where("team_id = ? AND user_id = ?", teamId, user.ID).
			Delete(&models.TeamMember{}).Error

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
	// Verify team exists
	var t models.Team
	err := db.First(&t, teamId).Error
	if err != nil {
		return err
	}

	// Verify user exists
	var u models.User
	err = db.First(&u, userId).Error
	if err != nil {
		return err
	}

	// Create TeamMember record with default member role
	teamMember := models.TeamMember{
		TeamID: teamId,
		UserID: userId,
		Role:   models.RoleMember,
	}

	err = db.Create(&teamMember).Error
	if err != nil {
		return err
	}

	return nil
}
