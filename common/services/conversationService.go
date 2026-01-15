package services

import (
	"fmt"
	"server/common/appError"
	"server/common/config"
	"server/common/models"
	"time"

	"gorm.io/gorm"
)

// CreateDirectConversation creates or returns existing direct conversation between two users
func CreateDirectConversation(currentUserID, otherUserID uint) (*models.Conversation, error) {
	if currentUserID == otherUserID {
		return nil, appError.ErrCannotMessageSelf
	}

	// Create a deterministic key for the direct conversation (smaller ID first)
	var userID1, userID2 uint
	if currentUserID < otherUserID {
		userID1, userID2 = currentUserID, otherUserID
	} else {
		userID1, userID2 = otherUserID, currentUserID
	}
	directKey := fmt.Sprintf("direct_%d_%d", userID1, userID2)

	var conversation models.Conversation

	// Check if conversation exists using count to avoid logging errors
	var count int64
	err := config.DB.Model(&models.Conversation{}).
		Where("direct_key = ?", directKey).
		Count(&count).Error
	if err != nil {
		return nil, err
	}

	// If conversation exists, load it with participants
	if count > 0 {
		err := config.DB.Where("direct_key = ?", directKey).
			Preload("Participants.User").
			First(&conversation).Error
		if err != nil {
			return nil, err
		}
		return &conversation, nil
	}

	// Create new conversation
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		conversation = models.Conversation{
			Type:      models.ConversationTypeDirect,
			DirectKey: &directKey,
		}

		if err := tx.Create(&conversation).Error; err != nil {
			return err
		}

		// Add both participants
		participants := []models.ConversationParticipant{
			{ConversationID: conversation.ID, UserID: currentUserID},
			{ConversationID: conversation.ID, UserID: otherUserID},
		}

		if err := tx.Create(&participants).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload with participants
	config.DB.Preload("Participants.User").First(&conversation, conversation.ID)

	return &conversation, nil
}

// CreateGroupConversation creates a new group conversation
func CreateGroupConversation(currentUserID uint, participantIDs []uint, title string) (*models.Conversation, error) {
	if len(participantIDs) < 1 {
		return nil, appError.ErrInsufficientParticipants
	}

	// Ensure current user is in the participant list
	hasCurrentUser := false
	for _, id := range participantIDs {
		if id == currentUserID {
			hasCurrentUser = true
			break
		}
	}
	if !hasCurrentUser {
		participantIDs = append(participantIDs, currentUserID)
	}

	var conversation models.Conversation

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		conversation = models.Conversation{
			Type:  models.ConversationTypeGroup,
			Title: &title,
		}

		if err := tx.Create(&conversation).Error; err != nil {
			return err
		}

		// Add all participants
		participants := make([]models.ConversationParticipant, len(participantIDs))
		for i, userID := range participantIDs {
			participants[i] = models.ConversationParticipant{
				ConversationID: conversation.ID,
				UserID:         userID,
			}
		}

		if err := tx.Create(&participants).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload with participants
	config.DB.Preload("Participants.User").First(&conversation, conversation.ID)

	return &conversation, nil
}

// EnsureTeamConversation ensures a team conversation exists for the given team
func EnsureTeamConversation(teamID uint) (*models.Conversation, error) {
	var conversation models.Conversation

	// Check if conversation exists using count to avoid logging errors
	var count int64
	err := config.DB.Model(&models.Conversation{}).
		Where("team_id = ?", teamID).
		Count(&count).Error
	if err != nil {
		return nil, err
	}

	// If conversation exists, load it with participants
	if count > 0 {
		err := config.DB.Where("team_id = ?", teamID).
			Preload("Participants.User").
			First(&conversation).Error
		if err != nil {
			return nil, err
		}
		return &conversation, nil
	}

	// Create new team conversation
	// Note: Participants will be synced separately via SyncTeamConversationMembers
	conversation = models.Conversation{
		Type:   models.ConversationTypeTeam,
		TeamID: &teamID,
	}

	if err := config.DB.Create(&conversation).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

// SyncTeamConversationMembers syncs team conversation participants with team members
func SyncTeamConversationMembers(teamID uint, memberIDs []uint) error {
	// Find or create team conversation
	conversation, err := EnsureTeamConversation(teamID)
	if err != nil {
		return err
	}

	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Get current participants
		var currentParticipants []models.ConversationParticipant
		if err := tx.Where("conversation_id = ?", conversation.ID).
			Find(&currentParticipants).Error; err != nil {
			return err
		}

		currentMemberMap := make(map[uint]bool)
		for _, p := range currentParticipants {
			currentMemberMap[p.UserID] = true
		}

		newMemberMap := make(map[uint]bool)
		for _, id := range memberIDs {
			newMemberMap[id] = true
		}

		// Add missing members
		for _, memberID := range memberIDs {
			if !currentMemberMap[memberID] {
				participant := models.ConversationParticipant{
					ConversationID: conversation.ID,
					UserID:         memberID,
				}
				if err := tx.Create(&participant).Error; err != nil {
					return err
				}
			}
		}

		// Remove members no longer in team
		for _, p := range currentParticipants {
			if !newMemberMap[p.UserID] {
				if err := tx.Delete(&p).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GetConversationByID retrieves a conversation by ID
func GetConversationByID(conversationID uint) (*models.Conversation, error) {
	var conversation models.Conversation
	err := config.DB.Preload("Participants.User").
		First(&conversation, conversationID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appError.ErrConversationNotFound
		}
		return nil, err
	}
	return &conversation, nil
}

// IsConversationMember checks if a user is a member of a conversation
func IsConversationMember(conversationID, userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ? AND left_at IS NULL", conversationID, userID).
		Count(&count).Error
	return count > 0, err
}

// ListConversations returns all conversations for a user with unread counts and last message
func ListConversations(userID uint) ([]models.Conversation, []int64, []*models.Message, error) {
	// Get all conversations where user is a participant AND conversation has at least one message
	var participants []models.ConversationParticipant
	err := config.DB.Where("user_id = ? AND left_at IS NULL", userID).
		// Only include conversations that have messages
		Joins("JOIN messages ON messages.conversation_id = conversation_participants.conversation_id").
		Preload("Conversation.Team").
		Preload("Conversation.Participants.User").
		Preload("Conversation.Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Order("conversation_id DESC").
		Distinct(). // Avoid duplicates from the JOIN
		Find(&participants).Error

	if err != nil {
		return nil, nil, nil, err
	}

	if len(participants) == 0 {
		return []models.Conversation{}, []int64{}, []*models.Message{}, nil
	}

	conversations := make([]models.Conversation, len(participants))
	unreadCounts := make([]int64, len(participants))
	lastMessages := make([]*models.Message, len(participants))

	for i, p := range participants {
		conversations[i] = p.Conversation

		// Calculate unread count
		query := config.DB.Model(&models.Message{}).
			Where("conversation_id = ? AND sender_id != ?", p.ConversationID, userID)

		if p.LastReadAt != nil {
			query = query.Where("created_at > ?", p.LastReadAt)
		} else {
			query = query.Where("created_at > ?", p.JoinedAt)
		}

		query.Count(&unreadCounts[i])

		// Get last message
		var lastMsg models.Message
		err := config.DB.Where("conversation_id = ?", p.ConversationID).
			Preload("Sender").
			Order("created_at DESC").
			First(&lastMsg).Error

		if err == nil {
			lastMessages[i] = &lastMsg
		}
	}

	return conversations, unreadCounts, lastMessages, nil
}

// MarkConversationRead updates the last_read_at timestamp for a user in a conversation
func MarkConversationRead(conversationID, userID uint, readAt time.Time) error {
	result := config.DB.Model(&models.ConversationParticipant{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("last_read_at", readAt)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return appError.ErrNotConversationMember
	}

	return nil
}

// GetConversationParticipantIDs returns all user IDs that currently belong to a conversation.
// - Only includes participants where left_at IS NULL
// - Used by the chat WebSocket hub for routing conversation messages live
func GetConversationParticipantIDs(conversationID uint) ([]uint, error) {
	var userIDs []uint

	err := config.DB.
		Model(&models.ConversationParticipant{}).
		Select("user_id").
		Where("conversation_id = ? AND left_at IS NULL", conversationID).
		Scan(&userIDs).Error

	if err != nil {
		return nil, err
	}

	return userIDs, nil
}
