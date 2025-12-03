package services

import (
	"server/common/appError"
	"server/common/config"
	"server/common/dto"
	"server/common/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateChat(creatorID uint, reqUserIDs []uint, name string) (*models.Chat, error) {
	var chat *models.Chat

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// Ensure creator is in the list
		userIDs := append([]uint{creatorID}, reqUserIDs...)

		// Remove duplicates and validate blocking
		uniqueIDs := make(map[uint]bool)
		var users []models.User
		for _, id := range userIDs {
			if !uniqueIDs[id] {
				// Check if the creator is blocked by this user (or vice versa)
				// Note: IsBlocked uses the global DB connection, which is acceptable for this check
				if id != creatorID && IsBlocked(id, creatorID) {
					continue // Skip blocked users
				}
				uniqueIDs[id] = true
				users = append(users, models.User{ID: id})
			}
		}

		newChat := models.Chat{
			Name:  name,
			Users: users,
		}

		if err := tx.Create(&newChat).Error; err != nil {
			return err
		}

		// Load full details for response using the same transaction to ensure consistency
		if err := tx.Preload("Users.FavoriteSports").First(&newChat, newChat.ID).Error; err != nil {
			return err
		}

		chat = &newChat
		return nil
	})

	if err != nil {
		return nil, err
	}

	return chat, nil
}

// GetUserChatsWithUnread returns chats with the unread count for the specific user
func GetUserChatsWithUnread(userID uint) ([]dto.ChatResponseDto, error) {
	var user models.User

	// 1. Fetch User with Chats and Preload Users
	// We verify the query logic: Get user -> Related Chats -> Related Users
	err := config.DB.
		Preload("Chats", func(db *gorm.DB) *gorm.DB {
			return db.Order("updated_at DESC")
		}).
		Preload("Chats.Users.FavoriteSports").
		First(&user, userID).Error

	if err != nil {
		return nil, err
	}

	var response []dto.ChatResponseDto

	for _, chat := range user.Chats {
		// 2. Get LastReadAt for this user/chat combo
		var userChat models.UserChat
		var lastReadAt time.Time

		// Check the join table for the specific LastReadAt timestamp
		err := config.DB.Where("chat_id = ? AND user_id = ?", chat.ID, userID).First(&userChat).Error
		if err == nil {
			lastReadAt = userChat.LastReadAt
		} else {
			// If record doesn't exist (legacy data), assume unread (time zero)
			lastReadAt = time.Time{}
		}

		// 3. Count messages created AFTER lastReadAt and NOT sent by the current user
		var count int64
		config.DB.Model(&models.Message{}).
			Where("chat_id = ? AND created_at > ? AND sender_id != ?", chat.ID, lastReadAt, userID).
			Count(&count)

		// 4. Convert to DTO
		response = append(response, dto.ToChatResponseDto(chat, count))
	}

	return response, nil
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

func MarkChatAsRead(chatID uint, userID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Upsert: Create the record if it doesn't exist, otherwise update the timestamp
		userChat := models.UserChat{
			ChatID:     chatID,
			UserID:     userID,
			LastReadAt: time.Now(),
		}

		// Use clause.OnConflict within the transaction
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "chat_id"}, {Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"last_read_at", "updated_at"}),
		}).Create(&userChat).Error
	})
}
