package integration

import (
	"testing"
	"time"

	"server/common/config"
	"server/common/models"
	"server/common/services"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage_PermissionCheck(t *testing.T) {
	setupTest(t)

	// Create test users
	user1 := models.User{Email: "user1@test.com", Password: "hash", FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: "hash", FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: "hash", FirstName: "User", LastName: "Three"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)
	config.DB.Create(&user3)

	// Create conversation between user1 and user2
	conv, _ := services.CreateDirectConversation(user1.ID, user2.ID)

	// User1 can send message
	msg1, err := services.SendMessage(conv.ID, user1.ID, "Hello from user1")
	assert.NoError(t, err)
	assert.NotNil(t, msg1)
	assert.Equal(t, "Hello from user1", msg1.Content)

	// User2 can send message
	msg2, err := services.SendMessage(conv.ID, user2.ID, "Hello from user2")
	assert.NoError(t, err)
	assert.NotNil(t, msg2)

	// User3 cannot send message (not a member)
	_, err = services.SendMessage(conv.ID, user3.ID, "Hello from user3")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a member")
}

func TestMarkConversationRead_UnreadCount(t *testing.T) {
	setupTest(t)

	// Create test users
	user1 := models.User{Email: "user1@test.com", Password: "hash", FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: "hash", FirstName: "User", LastName: "Two"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)

	// Create conversation
	conv, _ := services.CreateDirectConversation(user1.ID, user2.ID)

	// User1 sends 3 messages
	services.SendMessage(conv.ID, user1.ID, "Message 1")
	time.Sleep(10 * time.Millisecond)
	services.SendMessage(conv.ID, user1.ID, "Message 2")
	time.Sleep(10 * time.Millisecond)
	services.SendMessage(conv.ID, user1.ID, "Message 3")

	// Get user2's participant record
	var participant models.ConversationParticipant
	config.DB.Where("conversation_id = ? AND user_id = ?", conv.ID, user2.ID).First(&participant)

	// Calculate unread count for user2
	var unreadCount int64
	query := config.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_id != ?", conv.ID, user2.ID)

	if participant.LastReadAt != nil {
		query = query.Where("created_at > ?", participant.LastReadAt)
	} else {
		query = query.Where("created_at > ?", participant.JoinedAt)
	}

	query.Count(&unreadCount)
	assert.Equal(t, int64(3), unreadCount, "User2 should have 3 unread messages")

	// User2 marks conversation as read
	err := services.MarkConversationRead(conv.ID, user2.ID, time.Now())
	assert.NoError(t, err)

	// Recalculate unread count
	config.DB.Where("conversation_id = ? AND user_id = ?", conv.ID, user2.ID).First(&participant)
	query = config.DB.Model(&models.Message{}).
		Where("conversation_id = ? AND sender_id != ?", conv.ID, user2.ID)

	if participant.LastReadAt != nil {
		query = query.Where("created_at > ?", participant.LastReadAt)
	}

	query.Count(&unreadCount)
	assert.Equal(t, int64(0), unreadCount, "User2 should have 0 unread messages after marking as read")
}

func TestGetMessages_Pagination(t *testing.T) {
	setupTest(t)

	// Create test users
	user1 := models.User{Email: "user1@test.com", Password: "hash", FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: "hash", FirstName: "User", LastName: "Two"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)

	// Create conversation
	conv, _ := services.CreateDirectConversation(user1.ID, user2.ID)

	// Send 10 messages
	for i := 1; i <= 10; i++ {
		services.SendMessage(conv.ID, user1.ID, "Message "+string(rune(i)))
		time.Sleep(5 * time.Millisecond)
	}

	// Get first 5 messages
	messages, hasMore, total, err := services.GetMessages(conv.ID, user1.ID, 5, nil)
	assert.NoError(t, err)
	assert.Len(t, messages, 5, "Should return 5 messages")
	assert.True(t, hasMore, "Should have more messages")
	assert.Equal(t, int64(10), total, "Total should be 10")

	// Get next 5 messages using cursor
	lastMsgID := messages[0].ID
	messages2, hasMore2, _, err := services.GetMessages(conv.ID, user1.ID, 5, &lastMsgID)
	assert.NoError(t, err)
	assert.Len(t, messages2, 5, "Should return 5 more messages")
	assert.False(t, hasMore2, "Should not have more messages")
}

func TestGetMessages_NonMemberDenied(t *testing.T) {
	setupTest(t)

	// Create test users
	user1 := models.User{Email: "user1@test.com", Password: "hash", FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: "hash", FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: "hash", FirstName: "User", LastName: "Three"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)
	config.DB.Create(&user3)

	// Create conversation between user1 and user2
	conv, _ := services.CreateDirectConversation(user1.ID, user2.ID)

	// Send some messages
	services.SendMessage(conv.ID, user1.ID, "Private message")

	// User3 tries to read messages
	_, _, _, err := services.GetMessages(conv.ID, user3.ID, 10, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a member")
}
