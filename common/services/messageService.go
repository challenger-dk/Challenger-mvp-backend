package services

import (
	"log/slog"
	"server/common/appError"
	"server/common/config"
	"server/common/models"
	"strings"

	"gorm.io/gorm"
)

const maxPushBodyLength = 200

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

	// Send push notifications to recipients (fire-and-forget, errors logged)
	sendMessagePushNotifications(&message)

	return &message, nil
}

// sendMessagePushNotifications sends push notifications to all conversation recipients except the sender.
// Recipients who have blocked the sender or have no Expo token are skipped.
// Errors are logged but do not affect the caller.
func sendMessagePushNotifications(message *models.Message) {
	if message.ConversationID == nil {
		return
	}

	participantIDs, err := GetConversationParticipantIDs(*message.ConversationID)
	if err != nil {
		slog.Warn("Failed to get conversation participants for push",
			slog.Uint64("conversation_id", uint64(*message.ConversationID)),
			slog.Any("error", err),
		)
		return
	}

	// Collect recipient IDs (exclude sender, exclude users who blocked the sender)
	var recipientIDs []uint
	for _, userID := range participantIDs {
		if userID == message.SenderID {
			continue
		}
		if IsBlocked(userID, message.SenderID) {
			continue
		}
		recipientIDs = append(recipientIDs, userID)
	}

	if len(recipientIDs) == 0 {
		return
	}

	// Load recipients with Expo tokens
	var recipients []models.User
	if err := config.DB.Select("id", "expo_token", "first_name", "last_name").
		Where("id IN ?", recipientIDs).
		Where("expo_token != ''").
		Find(&recipients).Error; err != nil {
		slog.Warn("Failed to load recipients for push",
			slog.Uint64("conversation_id", uint64(*message.ConversationID)),
			slog.Any("error", err),
		)
		return
	}

	senderName := getSenderDisplayName(message)
	body := message.Content
	if len(body) > maxPushBodyLength {
		body = body[:maxPushBodyLength-3] + "..."
	}

	data := map[string]any{
		"conversation_id": *message.ConversationID,
		"message_id":      message.ID,
		"sender_id":       message.SenderID,
	}

	for _, recipient := range recipients {
		if recipient.ExpoToken == "" {
			continue
		}
		err := SendExpoPushNotification(
			recipient.ExpoToken,
			senderName,
			body,
			data,
		)
		if err != nil {
			slog.Warn("Failed to send message push notification",
				slog.Uint64("recipient_id", uint64(recipient.ID)),
				slog.String("expo_token_prefix", truncateForLog(recipient.ExpoToken, 30)),
				slog.Any("error", err),
			)
			if IsDeviceNotRegistered(err) {
				config.DB.Model(&models.User{}).Where("id = ?", recipient.ID).Update("expo_token", "")
			}
		}
	}
}

func getSenderDisplayName(message *models.Message) string {
	if message.Sender.ID != 0 {
		name := strings.TrimSpace(message.Sender.FirstName + " " + message.Sender.LastName)
		if name != "" {
			return name
		}
		if message.Sender.FirstName != "" {
			return message.Sender.FirstName
		}
	}
	return "Ny besked"
}

func truncateForLog(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
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

	// Get total count (exclude messages from blocked users)
	var total int64
	config.DB.Model(&models.Message{}).
		Scopes(ExcludeBlockedUsersOn(userID, "sender_id")).
		Where("conversation_id = ?", conversationID).
		Count(&total)

	// Build query (exclude messages from blocked users)
	query := config.DB.
		Scopes(ExcludeBlockedUsersOn(userID, "sender_id")).
		Where("conversation_id = ?", conversationID).
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
