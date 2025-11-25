package services

import (
	"server/common/appError"
	"server/common/config"
	"server/common/dto"
	"server/common/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsers() ([]models.User, error) {
	var users []models.User

	err := config.DB.Preload("FavoriteSports").
		Find(&users).
		Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(userID uint) (*models.User, error) {
	var user models.User

	err := config.DB.Preload("FavoriteSports").
		Preload("Friends").
		First(&user, userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByIDWithSettings(userID uint) (*models.User, error) {
	var user models.User

	err := config.DB.Preload("FavoriteSports").
		Preload("Friends").
		Preload("Settings").
		First(&user, userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserSettings(userID uint) (*models.UserSettings, error) {
	var settings models.UserSettings
	err := config.DB.First(&settings, userID).Error

	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func GetInCommonStats(currentUserID, targetUserID uint) (dto.CommonStatsDto, error) {
	var stats dto.CommonStatsDto
	db := config.DB

	// 1. Common Teams Count
	// Corrected table name from "team_users" to "user_teams"
	var count int64
	err := db.Table("user_teams as t1").
		Joins("JOIN user_teams as t2 ON t1.team_id = t2.team_id").
		Where("t1.user_id = ? AND t2.user_id = ?", currentUserID, targetUserID).
		Count(&count).Error

	if err != nil {
		return stats, err
	}
	stats.CommonTeamsCount = count

	// 2. Common Friends Count
	err = db.Table("user_friends as f1").
		Joins("JOIN user_friends as f2 ON f1.friend_id = f2.friend_id").
		Where("f1.user_id = ? AND f2.user_id = ?", currentUserID, targetUserID).
		Count(&count).Error

	if err != nil {
		return stats, err
	}
	stats.CommonFriendsCount = count

	// 3. Common Sports (Favorites)
	var commonSports []models.Sport
	err = db.Table("sports").
		Joins("JOIN user_favorite_sports as us1 ON us1.sport_id = sports.id").
		Joins("JOIN user_favorite_sports as us2 ON us2.sport_id = sports.id").
		Where("us1.user_id = ? AND us2.user_id = ?", currentUserID, targetUserID).
		Find(&commonSports).Error

	if err != nil {
		return stats, err
	}

	// Convert models to DTOs
	stats.CommonSports = make([]dto.SportDto, len(commonSports))
	for i, s := range commonSports {
		stats.CommonSports[i] = dto.ToSportDto(s)
	}

	return stats, nil
}

func CreateUser(email, password, firstName, lastName string, favoriteSports []string) (*models.User, error) {
	var user *models.User

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var existingUser models.User

		err := tx.Where("email = ?", email).
			First(&existingUser).
			Error

		if err == nil {
			return appError.ErrUserExists
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		newUser := &models.User{
			Email:     email,
			Password:  string(hashedPassword),
			FirstName: firstName,
			LastName:  lastName,
			Settings:  &models.UserSettings{},
		}

		err = tx.Create(newUser).Error
		if err != nil {
			return err
		}

		// Associate favorite sports if provided
		if len(favoriteSports) > 0 {
			if err := associateFavoriteSports(tx, newUser.ID, favoriteSports); err != nil {
				return err
			}
			// Reload user with favorite sports
			if err := tx.Preload("FavoriteSports").First(newUser, newUser.ID).Error; err != nil {
				return err
			}
		}

		user = newUser
		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(userID uint, user dto.UserUpdateDto) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var existingUser models.User
		if err := tx.Preload("FavoriteSports").First(&existingUser, userID).Error; err != nil {
			return err
		}

		if user.FirstName != "" {
			existingUser.FirstName = user.FirstName
		}

		if user.LastName != "" {
			existingUser.LastName = user.LastName
		}

		if user.ProfilePicture != "" {
			existingUser.ProfilePicture = user.ProfilePicture
		}

		if user.Bio != "" {
			existingUser.Bio = user.Bio
		}

		if err := tx.Save(&existingUser).Error; err != nil {
			return err
		}

		// Update favorite sports if provided
		if user.FavoriteSports != nil {
			if err := associateFavoriteSports(tx, userID, user.FavoriteSports); err != nil {
				return err
			}
		}

		return nil
	})
}

func UpdateUserSettings(userID uint, settingsDto dto.UserSettingsUpdateDto) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var settings models.UserSettings

		if err := tx.First(&settings, userID).Error; err != nil {
			return err
		}

		if settingsDto.NotifyTeamInvite != nil {
			settings.NotifyTeamInvite = *settingsDto.NotifyTeamInvite
		}
		if settingsDto.NotifyFriendReq != nil {
			settings.NotifyFriendReq = *settingsDto.NotifyFriendReq
		}
		if settingsDto.NotifyChallengeInvite != nil {
			settings.NotifyChallengeInvite = *settingsDto.NotifyChallengeInvite
		}
		if settingsDto.NotifyChallengeUpdate != nil {
			settings.NotifyChallengeUpdate = *settingsDto.NotifyChallengeUpdate
		}

		return tx.Save(&settings).Error
	})
}

func DeleteUser(userID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}

		// Clear all many-to-many associations
		if err := tx.Model(&user).Association("FavoriteSports").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&user).Association("Friends").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&user).Association("Teams").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&user).Association("JoinedChallenges").Clear(); err != nil {
			return err
		}

		// Handle one-to-many relationships (delete or set to null)
		if err := tx.Where("creator_id = ?", userID).Delete(&models.Challenge{}).Error; err != nil {
			return err
		}
		if err := tx.Where("creator_id = ?", userID).Delete(&models.Team{}).Error; err != nil {
			return err
		}

		// Clean up invitations
		if err := tx.Where("inviter_id = ? OR invitee_id = ?", userID, userID).Delete(&models.Invitation{}).Error; err != nil {
			return err
		}

		// Delete the user
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}

		return nil
	})
}

// DeleteFriendship removes both users from each other's friends list
func RemoveFriend(userIdA uint, userIdB uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {

		// Ids must be different
		if userIdA == userIdB {
			return appError.ErrInvalidFriendship
		}

		var userA, userB models.User

		if err := tx.First(&userA, userIdA).Error; err != nil {
			return err
		}

		if err := tx.First(&userB, userIdB).Error; err != nil {
			return err
		}

		err := tx.Model(&userA).
			Association("Friends").
			Delete(&userB)

		if err != nil {
			return err
		}

		err = tx.Model(&userB).
			Association("Friends").
			Delete(&userA)

		if err != nil {
			return err
		}

		return nil
	})
}

// Package private
// createFriendship adds both users to each other's friends list
func createFriendship(userIdA uint, userIdB uint, db *gorm.DB) error {
	var userA, userB models.User

	if err := db.First(&userA, userIdA).Error; err != nil {
		return err
	}
	if err := db.First(&userB, userIdB).Error; err != nil {
		return err
	}

	err := db.Model(&userA).
		Association("Friends").
		Append(&userB)

	if err != nil {
		return err
	}

	err = db.Model(&userB).
		Association("Friends").
		Append(&userA)

	if err != nil {
		return err
	}

	return nil
}
