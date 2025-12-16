package services

import (
	"server/common/appError"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

// SendMessage creates a new message in a conversation
func SendMessage(conversationID, senderID uint, content string) (*models.Message, error) {
	// Check if sender is a member of the conversation
	isMember, err := IsConversationMember(conversationID, senderID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, appError.ErrNotConversationMember
	}

	message := models.Message{
		ConversationID: &conversationID,
		SenderID:       senderID,
		Content:        content,
	}

	if err := config.DB.Create(&message).Error; err != nil {
		return nil, err
	}

	// Preload sender info
	config.DB.Preload("Sender").First(&message, message.ID)

	// Update conversation's updated_at timestamp
	config.DB.Model(&models.Conversation{}).
		Where("id = ?", conversationID).
		Update("updated_at", message.CreatedAt)

	return &message, nil
}

// GetMessages retrieves messages from a conversation with pagination
func GetMessages(conversationID, userID uint, limit int, beforeMessageID *uint) ([]models.Message, bool, int64, error) {
	// Check if user is a member
	isMember, err := IsConversationMember(conversationID, userID)
	if err != nil {
		return nil, false, 0, err
	}
	if !isMember {
		return nil, false, 0, appError.ErrNotConversationMember
	}

	// Get total count
	var total int64
	config.DB.Model(&models.Message{}).
		Where("conversation_id = ?", conversationID).
		Count(&total)

	// Build query
	query := config.DB.Where("conversation_id = ?", conversationID).
		Preload("Sender").
		Order("created_at DESC")

	// Apply cursor-based pagination
	if beforeMessageID != nil {
		var beforeMsg models.Message
		if err := config.DB.First(&beforeMsg, *beforeMessageID).Error; err == nil {
			query = query.Where("created_at < ?", beforeMsg.CreatedAt)
		}
	}

	// Fetch limit + 1 to check if there are more
	query = query.Limit(limit + 1)

	var messages []models.Message
	if err := query.Find(&messages).Error; err != nil {
		return nil, false, 0, err
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	// Reverse to get chronological order (oldest first)
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, hasMore, total, nil
}

// GetMessageByID retrieves a single message
func GetMessageByID(messageID uint) (*models.Message, error) {
	var message models.Message
	err := config.DB.Preload("Sender").First(&message, messageID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appError.ErrConversationNotFound // Reusing error
		}
		return nil, err
	}
	return &message, nil
}

