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
	user1 := models.User{Email: "user1@test.com", Password: "hash", FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: "hash", FirstName: "User", LastName: "Two"}
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

	user := models.User{Email: "user@test.com", Password: "hash", FirstName: "User", LastName: "One"}
	config.DB.Create(&user)

	_, err := services.CreateDirectConversation(user.ID, user.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create conversation with yourself")
}

func TestIsConversationMember(t *testing.T) {
	setupTest(t)

	// Create test users
	user1 := models.User{Email: "user1@test.com", Password: "hash", FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: "hash", FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: "hash", FirstName: "User", LastName: "Three"}
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
	user1 := models.User{Email: "user1@test.com", Password: "hash", FirstName: "User", LastName: "One"}
	user2 := models.User{Email: "user2@test.com", Password: "hash", FirstName: "User", LastName: "Two"}
	user3 := models.User{Email: "user3@test.com", Password: "hash", FirstName: "User", LastName: "Three"}
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
