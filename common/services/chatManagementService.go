package services

import (
	"server/common/appError"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

func CreateChat(creatorID uint, reqUserIDs []uint, name string) (*models.Chat, error) {
	// Ensure creator is in the list
	userIDs := append([]uint{creatorID}, reqUserIDs...)

	// Remove duplicates and validate blocking
	uniqueIDs := make(map[uint]bool)
	var users []models.User
	for _, id := range userIDs {
		if !uniqueIDs[id] {
			// Check if the creator is blocked by this user (or vice versa)
			if id != creatorID && IsBlocked(id, creatorID) {
				continue // Skip blocked users
			}
			uniqueIDs[id] = true
			users = append(users, models.User{ID: id})
		}
	}

	chat := models.Chat{
		Name:  name,
		Users: users,
	}

	if err := config.DB.Create(&chat).Error; err != nil {
		return nil, err
	}

	// Load full details for response
	if err := config.DB.Preload("Users.FavoriteSports").First(&chat, chat.ID).Error; err != nil {
		return nil, err
	}

	return &chat, nil
}

func GetUserChats(userID uint) ([]models.Chat, error) {
	var user models.User
	err := config.DB.Preload("Chats.Users.FavoriteSports").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.Chats, nil
}

func GetChatByID(chatID uint, userID uint) (*models.Chat, error) {
	var chat models.Chat
	err := config.DB.Preload("Users.FavoriteSports").First(&chat, chatID).Error
	if err != nil {
		return nil, err
	}

	// Authorization check
	isMember := false
	for _, u := range chat.Users {
		if u.ID == userID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, appError.ErrUnauthorized
	}

	return &chat, nil
}

func AddUserToChat(chatID uint, currentUserID uint, newUserID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		var chat models.Chat
		if err := tx.Preload("Users").First(&chat, chatID).Error; err != nil {
			return err
		}

		// Auth check: Is requester a member?
		isMember := false
		for _, u := range chat.Users {
			if u.ID == currentUserID {
				isMember = true
				break
			}
		}
		if !isMember {
			return appError.ErrUnauthorized
		}

		// Check blocking
		if IsBlocked(newUserID, currentUserID) {
			return appError.ErrUserBlocked
		}

		// Add new user
		var newUser models.User
		if err := tx.First(&newUser, newUserID).Error; err != nil {
			return err
		}

		return tx.Model(&chat).Association("Users").Append(&newUser)
	})
}
