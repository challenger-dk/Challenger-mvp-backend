package integration

import (
	"testing"

	"server/common/config"
	"server/common/models"
	"server/common/services"

	"github.com/stretchr/testify/assert"
)

func TestCreateDirectConversation_Idempotent(t *testing.T) {
	setupTest(t)

	// Create test users
	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: &password, FirstName: "User", LastName: "Two"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)

	// Create conversation first time
	conv1, err := services.CreateDirectConversation(user1.ID, user2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, conv1)
	assert.Equal(t, models.ConversationTypeDirect, conv1.Type)

	// Create conversation second time - should return same conversation
	conv2, err := services.CreateDirectConversation(user1.ID, user2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, conv2)
	assert.Equal(t, conv1.ID, conv2.ID, "Should return same conversation ID")

	// Create from opposite direction - should still return same conversation
	conv3, err := services.CreateDirectConversation(user2.ID, user1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, conv3)
	assert.Equal(t, conv1.ID, conv3.ID, "Should return same conversation ID regardless of order")
}

func TestCreateDirectConversation_CannotMessageSelf(t *testing.T) {
	setupTest(t)

	password := "hash"
	user := models.User{Email: "user@test.com", Password: &password, FirstName: "User", LastName: "One"}
	config.DB.Create(&user)

	_, err := services.CreateDirectConversation(user.ID, user.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create conversation with yourself")
}

func TestIsConversationMember(t *testing.T) {
	setupTest(t)

	// Create test users
	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: &password, FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: &password, FirstName: "User", LastName: "Three"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)
	config.DB.Create(&user3)

	// Create conversation
	conv, _ := services.CreateDirectConversation(user1.ID, user2.ID)

	// Check membership
	isMember1, err := services.IsConversationMember(conv.ID, user1.ID)
	assert.NoError(t, err)
	assert.True(t, isMember1, "User1 should be a member")

	isMember2, err := services.IsConversationMember(conv.ID, user2.ID)
	assert.NoError(t, err)
	assert.True(t, isMember2, "User2 should be a member")

	isMember3, err := services.IsConversationMember(conv.ID, user3.ID)
	assert.NoError(t, err)
	assert.False(t, isMember3, "User3 should not be a member")
}

func TestSyncTeamConversationMembers(t *testing.T) {
	setupTest(t)

	// Create test users
	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: &password, FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: &password, FirstName: "User", LastName: "Three"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)
	config.DB.Create(&user3)

	// Create team
	team := models.Team{Name: "Test Team", CreatorID: user1.ID}
	config.DB.Create(&team)

	// Sync with initial members
	err := services.SyncTeamConversationMembers(team.ID, []uint{user1.ID, user2.ID})
	assert.NoError(t, err)

	// Verify conversation was created
	var conv models.Conversation
	err = config.DB.Where("team_id = ?", team.ID).First(&conv).Error
	assert.NoError(t, err)

	// Verify participants
	var participants []models.ConversationParticipant
	config.DB.Where("conversation_id = ?", conv.ID).Find(&participants)
	assert.Len(t, participants, 2, "Should have 2 participants")

	// Sync with updated members (add user3, remove user2)
	err = services.SyncTeamConversationMembers(team.ID, []uint{user1.ID, user3.ID})
	assert.NoError(t, err)

	// Verify updated participants
	config.DB.Where("conversation_id = ?", conv.ID).Find(&participants)
	assert.Len(t, participants, 2, "Should still have 2 participants")

	// Check specific members
	isMember1, _ := services.IsConversationMember(conv.ID, user1.ID)
	isMember2, _ := services.IsConversationMember(conv.ID, user2.ID)
	isMember3, _ := services.IsConversationMember(conv.ID, user3.ID)

	assert.True(t, isMember1, "User1 should still be a member")
	assert.False(t, isMember2, "User2 should be removed")
	assert.True(t, isMember3, "User3 should be added")
}

func TestCreateGroupConversation_Success(t *testing.T) {
	setupTest(t)

	// Create test users
	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: &password, FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: &password, FirstName: "User", LastName: "Three"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)
	config.DB.Create(&user3)

	// Create group conversation
	title := "Test Group"
	conv, err := services.CreateGroupConversation(user1.ID, []uint{user2.ID, user3.ID}, title)
	assert.NoError(t, err)
	assert.NotNil(t, conv)
	assert.Equal(t, models.ConversationTypeGroup, conv.Type)
	assert.NotNil(t, conv.Title)
	assert.Equal(t, title, *conv.Title)

	// Verify participants are preloaded
	assert.NotNil(t, conv.Participants)
	assert.Len(t, conv.Participants, 3, "Should have 3 participants (user1, user2, user3)")

	// Verify all users are participants
	participantIDs := make(map[uint]bool)
	for _, p := range conv.Participants {
		participantIDs[p.UserID] = true
		assert.NotNil(t, p.User, "User should be preloaded")
	}
	assert.True(t, participantIDs[user1.ID], "User1 should be a participant")
	assert.True(t, participantIDs[user2.ID], "User2 should be a participant")
	assert.True(t, participantIDs[user3.ID], "User3 should be a participant")
}

func TestCreateGroupConversation_CurrentUserAutoAdded(t *testing.T) {
	setupTest(t)

	// Create test users
	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: &password, FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: &password, FirstName: "User", LastName: "Three"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)
	config.DB.Create(&user3)

	// Create group conversation without including user1 in participant list
	title := "Test Group"
	conv, err := services.CreateGroupConversation(user1.ID, []uint{user2.ID, user3.ID}, title)
	assert.NoError(t, err)
	assert.NotNil(t, conv)

	// Verify user1 is automatically added
	var participants []models.ConversationParticipant
	config.DB.Where("conversation_id = ?", conv.ID).Find(&participants)
	assert.Len(t, participants, 3, "Should have 3 participants")

	participantIDs := make(map[uint]bool)
	for _, p := range participants {
		participantIDs[p.UserID] = true
	}
	assert.True(t, participantIDs[user1.ID], "User1 should be automatically added")
}

func TestCreateGroupConversation_InsufficientParticipants(t *testing.T) {
	setupTest(t)

	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	config.DB.Create(&user1)

	// Try to create group conversation with empty participant list
	_, err := services.CreateGroupConversation(user1.ID, []uint{}, "Test Group")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "group conversation requires at least 2 participants")
}

func TestGetConversationByID_Success(t *testing.T) {
	setupTest(t)

	// Create test users
	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: &password, FirstName: "User", LastName: "Two"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)

	// Create direct conversation
	conv, err := services.CreateDirectConversation(user1.ID, user2.ID)
	assert.NoError(t, err)
	assert.NotNil(t, conv)

	// Get conversation by ID
	retrievedConv, err := services.GetConversationByID(conv.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedConv)
	assert.Equal(t, conv.ID, retrievedConv.ID)
	assert.Equal(t, models.ConversationTypeDirect, retrievedConv.Type)

	// Verify participants are preloaded
	assert.NotNil(t, retrievedConv.Participants)
	assert.Len(t, retrievedConv.Participants, 2, "Should have 2 participants")
	for _, p := range retrievedConv.Participants {
		assert.NotNil(t, p.User, "User should be preloaded")
	}
}

func TestGetConversationByID_GroupConversation(t *testing.T) {
	setupTest(t)

	// Create test users
	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: &password, FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: &password, FirstName: "User", LastName: "Three"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)
	config.DB.Create(&user3)

	// Create group conversation
	title := "Test Group"
	conv, err := services.CreateGroupConversation(user1.ID, []uint{user2.ID, user3.ID}, title)
	assert.NoError(t, err)
	assert.NotNil(t, conv)

	// Get conversation by ID
	retrievedConv, err := services.GetConversationByID(conv.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedConv)
	assert.Equal(t, conv.ID, retrievedConv.ID)
	assert.Equal(t, models.ConversationTypeGroup, retrievedConv.Type)
	assert.NotNil(t, retrievedConv.Title)
	assert.Equal(t, title, *retrievedConv.Title)

	// Verify participants are preloaded
	assert.NotNil(t, retrievedConv.Participants)
	assert.Len(t, retrievedConv.Participants, 3, "Should have 3 participants")
}

func TestGetConversationByID_NotFound(t *testing.T) {
	setupTest(t)

	// Try to get non-existent conversation
	_, err := services.GetConversationByID(99999)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conversation not found")
}

func TestListConversations_EmptyList(t *testing.T) {
	setupTest(t)

	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	config.DB.Create(&user1)

	// List conversations for user with no conversations
	conversations, unreadCounts, lastMessages, err := services.ListConversations(user1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, conversations)
	assert.NotNil(t, unreadCounts)
	assert.NotNil(t, lastMessages)
	assert.Len(t, conversations, 0, "Should return empty list")
	assert.Len(t, unreadCounts, 0, "Should return empty unread counts")
	assert.Len(t, lastMessages, 0, "Should return empty last messages")
}

func TestListConversations_WithMessages(t *testing.T) {
	setupTest(t)

	// Create test users
	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: &password, FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: &password, FirstName: "User", LastName: "Three"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)
	config.DB.Create(&user3)

	// Create two conversations
	conv1, _ := services.CreateDirectConversation(user1.ID, user2.ID)
	conv2, _ := services.CreateDirectConversation(user1.ID, user3.ID)

	// Send messages to both conversations
	services.SendMessage(conv1.ID, user1.ID, "Message 1 in conv1")
	services.SendMessage(conv1.ID, user2.ID, "Message 2 in conv1")
	services.SendMessage(conv2.ID, user1.ID, "Message 1 in conv2")
	services.SendMessage(conv2.ID, user3.ID, "Message 2 in conv2")

	// List conversations for user1
	conversations, unreadCounts, lastMessages, err := services.ListConversations(user1.ID)
	assert.NoError(t, err)
	assert.Len(t, conversations, 2, "Should return 2 conversations")

	// Verify unread counts and last messages are returned
	assert.Len(t, unreadCounts, 2, "Should return 2 unread counts")
	assert.Len(t, lastMessages, 2, "Should return 2 last messages")

	// Verify conversations have participants preloaded
	for _, conv := range conversations {
		assert.NotNil(t, conv.Participants, "Participants should be preloaded")
		for _, p := range conv.Participants {
			assert.NotNil(t, p.User, "User should be preloaded")
		}
	}

	// Verify last messages are populated
	for _, msg := range lastMessages {
		assert.NotNil(t, msg, "Last message should not be nil")
		assert.NotNil(t, msg.Sender, "Sender should be preloaded")
	}
}

func TestListConversations_UnreadCounts(t *testing.T) {
	setupTest(t)

	// Create test users
	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: &password, FirstName: "User", LastName: "Two"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)

	// Create conversation
	conv, _ := services.CreateDirectConversation(user1.ID, user2.ID)

	// User2 sends 3 messages (user1 hasn't read them)
	services.SendMessage(conv.ID, user2.ID, "Unread message 1")
	services.SendMessage(conv.ID, user2.ID, "Unread message 2")
	services.SendMessage(conv.ID, user2.ID, "Unread message 3")

	// List conversations for user1
	conversations, unreadCounts, _, err := services.ListConversations(user1.ID)
	assert.NoError(t, err)
	assert.Len(t, conversations, 1, "Should return 1 conversation")
	assert.Len(t, unreadCounts, 1, "Should return 1 unread count")
	assert.Equal(t, int64(3), unreadCounts[0], "User1 should have 3 unread messages")

	// User1 sends a message (shouldn't count as unread for user1)
	services.SendMessage(conv.ID, user1.ID, "Read message from user1")

	// Re-list conversations - unread count should still be 3
	conversations, unreadCounts, _, err = services.ListConversations(user1.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), unreadCounts[0], "User1 should still have 3 unread messages (own message doesn't count)")
}

func TestListConversations_OnlyConversationsWithMessages(t *testing.T) {
	setupTest(t)

	// Create test users
	password := "hash"
	user1 := models.User{Email: "user1@test.com", Password: &password, FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: &password, FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: &password, FirstName: "User", LastName: "Three"}
	config.DB.Create(&user1)
	config.DB.Create(&user2)
	config.DB.Create(&user3)

	// Create two conversations
	conv1, _ := services.CreateDirectConversation(user1.ID, user2.ID)
	_, _ = services.CreateDirectConversation(user1.ID, user3.ID)

	// Only send message to conv1
	services.SendMessage(conv1.ID, user1.ID, "Message in conv1")

	// List conversations for user1 - should only return conv1 (has messages)
	conversations, unreadCounts, lastMessages, err := services.ListConversations(user1.ID)
	assert.NoError(t, err)
	assert.Len(t, conversations, 1, "Should only return conversation with messages")
	assert.Equal(t, conv1.ID, conversations[0].ID, "Should return conv1")
	assert.Len(t, unreadCounts, 1, "Should return 1 unread count")
	assert.Len(t, lastMessages, 1, "Should return 1 last message")
}
