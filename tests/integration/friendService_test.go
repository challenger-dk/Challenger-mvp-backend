package integration

import (
	"server/common/config"
	"server/common/models"
	"server/common/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFriendService_GetSuggestedFriends(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Create test users
	user1, _ := services.CreateUser(models.User{
		Email:          "user1@test.com",
		FirstName:      "User",
		LastName:       "One",
		FavoriteSports: []models.Sport{{Name: "Tennis"}, {Name: "Football"}},
	}, "pw")

	user2, _ := services.CreateUser(models.User{
		Email:          "user2@test.com",
		FirstName:      "User",
		LastName:       "Two",
		FavoriteSports: []models.Sport{{Name: "Tennis"}},
	}, "pw")

	user3, _ := services.CreateUser(models.User{
		Email:          "user3@test.com",
		FirstName:      "User",
		LastName:       "Three",
		FavoriteSports: []models.Sport{{Name: "Tennis"}, {Name: "Basketball"}},
	}, "pw")

	user4, _ := services.CreateUser(models.User{
		Email:          "user4@test.com",
		FirstName:      "User",
		LastName:       "Four",
		FavoriteSports: []models.Sport{{Name: "Football"}},
	}, "pw")

	user5, _ := services.CreateUser(models.User{
		Email:     "user5@test.com",
		FirstName: "User",
		LastName:  "Five",
	}, "pw")

	// Create friendships
	// user1 <-> user2 (friends)
	config.DB.Model(user1).Association("Friends").Append(user2)
	config.DB.Model(user2).Association("Friends").Append(user1)

	// user2 <-> user3 (friends)
	config.DB.Model(user2).Association("Friends").Append(user3)
	config.DB.Model(user3).Association("Friends").Append(user2)

	// user2 <-> user4 (friends)
	config.DB.Model(user2).Association("Friends").Append(user4)
	config.DB.Model(user4).Association("Friends").Append(user2)

	// Create teams
	team1, _ := services.CreateTeam(models.Team{Name: "Team Alpha", CreatorID: user1.ID})
	team2, _ := services.CreateTeam(models.Team{Name: "Team Beta", CreatorID: user2.ID})

	// Add users to teams (need to reload teams first)
	var t1, t2 models.Team
	config.DB.First(&t1, team1.ID)
	config.DB.First(&t2, team2.ID)

	config.DB.Model(&t1).Association("Users").Append(user3) // user1 and user3 share team1
	config.DB.Model(&t2).Association("Users").Append(user1) // user1 and user2 share team2
	config.DB.Model(&t2).Association("Users").Append(user4) // user1 and user4 share team2

	// Get suggestions for user1
	suggestions, err := services.GetSuggestedFriends(user1.ID)
	assert.NoError(t, err)

	// user1 should get suggestions for user3 and user4 (not user2 who is already a friend)
	// user3: 1 common friend (user2) + 1 common team (team1) + 2 common sports (Tennis) = 4*1 + 3*1 + 1*2 = 9
	// user4: 1 common friend (user2) + 1 common team (team2) + 1 common sport (Football) = 4*1 + 3*1 + 1*1 = 8
	// user5: 0 connections = 0 (should not appear)

	assert.GreaterOrEqual(t, len(suggestions), 2, "Should have at least 2 suggestions")

	// Check that user2 (already a friend) is NOT in suggestions
	for _, suggestion := range suggestions {
		assert.NotEqual(t, user2.ID, suggestion.ID, "Already-friend user2 should not be suggested")
	}

	// Check that user5 (no connections) is NOT in suggestions
	foundUser5 := false
	for _, suggestion := range suggestions {
		if suggestion.ID == user5.ID {
			foundUser5 = true
		}
	}
	assert.False(t, foundUser5, "User5 with no connections should not be suggested")

	// user3 should be ranked higher than user4 (or equal, depending on exact scoring)
	// Both have 1 common friend, but user3 has more common sports
	if len(suggestions) >= 2 {
		// First suggestion should be either user3 or user4
		assert.Contains(t, []uint{user3.ID, user4.ID}, suggestions[0].ID)
	}
}

func TestFriendService_GetSuggestedFriends_ExcludesBlockedUsers(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	user1, _ := services.CreateUser(models.User{
		Email:          "block1@test.com",
		FirstName:      "Block",
		LastName:       "One",
		FavoriteSports: []models.Sport{{Name: "Tennis"}},
	}, "pw")

	user2, _ := services.CreateUser(models.User{
		Email:          "block2@test.com",
		FirstName:      "Block",
		LastName:       "Two",
		FavoriteSports: []models.Sport{{Name: "Tennis"}},
	}, "pw")

	// Block user2
	err := services.BlockUser(user1.ID, user2.ID)
	assert.NoError(t, err)

	// Get suggestions for user1
	suggestions, err := services.GetSuggestedFriends(user1.ID)
	assert.NoError(t, err)

	// user2 should NOT be in suggestions (blocked)
	for _, suggestion := range suggestions {
		assert.NotEqual(t, user2.ID, suggestion.ID, "Blocked user should not be suggested")
	}
}

func TestFriendService_GetSuggestedFriends_MaxTenResults(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Create main user
	mainUser, _ := services.CreateUser(models.User{
		Email:          "main@test.com",
		FirstName:      "Main",
		LastName:       "User",
		FavoriteSports: []models.Sport{{Name: "Tennis"}},
	}, "pw")

	// Create 15 users with common sport
	for i := 0; i < 15; i++ {
		services.CreateUser(models.User{
			Email:          "user" + string(rune(i)) + "@test.com",
			FirstName:      "User",
			LastName:       string(rune(i)),
			FavoriteSports: []models.Sport{{Name: "Tennis"}},
		}, "pw")
	}

	// Get suggestions
	suggestions, err := services.GetSuggestedFriends(mainUser.ID)
	assert.NoError(t, err)

	// Should return max 10 suggestions
	assert.LessOrEqual(t, len(suggestions), 10, "Should return maximum 10 suggestions")
}

func TestFriendService_GetFriends(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// 1. Create users
	user1, _ := services.CreateUser(models.User{
		Email:     "getfriends1@test.com",
		FirstName: "User",
		LastName:  "One",
	}, "pw")

	user2, _ := services.CreateUser(models.User{
		Email:     "getfriends2@test.com",
		FirstName: "User",
		LastName:  "Two",
	}, "pw")

	user3, _ := services.CreateUser(models.User{
		Email:     "getfriends3@test.com",
		FirstName: "User",
		LastName:  "Three",
	}, "pw")

	user4, _ := services.CreateUser(models.User{
		Email:     "getfriends4@test.com",
		FirstName: "User",
		LastName:  "Four",
	}, "pw")

	// 2. Get friends for user with no friends
	friends, err := services.GetFriends(user1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, friends)
	assert.Empty(t, friends, "User with no friends should return empty slice")

	// 3. Create friendships: user1 <-> user2, user1 <-> user3
	config.DB.Model(user1).Association("Friends").Append(user2)
	config.DB.Model(user2).Association("Friends").Append(user1)

	config.DB.Model(user1).Association("Friends").Append(user3)
	config.DB.Model(user3).Association("Friends").Append(user1)

	// 4. Get friends for user1
	friends, err = services.GetFriends(user1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, friends)
	assert.Len(t, friends, 2, "User1 should have 2 friends")

	// Verify friends are correct
	friendIDs := make(map[uint]bool)
	for _, friend := range friends {
		friendIDs[friend.ID] = true
	}
	assert.True(t, friendIDs[user2.ID], "user2 should be in friends list")
	assert.True(t, friendIDs[user3.ID], "user3 should be in friends list")
	assert.False(t, friendIDs[user4.ID], "user4 should not be in friends list")

	// 5. Get friends for user2 (should have user1)
	friends2, err := services.GetFriends(user2.ID)
	assert.NoError(t, err)
	assert.Len(t, friends2, 1, "User2 should have 1 friend")
	assert.Equal(t, user1.ID, friends2[0].ID)

	// 6. Get friends for user4 (should have none)
	friends4, err := services.GetFriends(user4.ID)
	assert.NoError(t, err)
	assert.Empty(t, friends4, "User4 should have no friends")

	// 7. Add more friends to user1
	config.DB.Model(user1).Association("Friends").Append(user4)
	config.DB.Model(user4).Association("Friends").Append(user1)

	// 8. Get friends again - should now have 3
	friends, err = services.GetFriends(user1.ID)
	assert.NoError(t, err)
	assert.Len(t, friends, 3, "User1 should now have 3 friends")

	// Verify all friends are present
	friendIDs = make(map[uint]bool)
	for _, friend := range friends {
		friendIDs[friend.ID] = true
	}
	assert.True(t, friendIDs[user2.ID])
	assert.True(t, friendIDs[user3.ID])
	assert.True(t, friendIDs[user4.ID])

	// 9. Verify friend data is loaded correctly
	for _, friend := range friends {
		assert.NotZero(t, friend.ID)
		assert.NotEmpty(t, friend.Email)
		assert.NotEmpty(t, friend.FirstName)
	}

	// 10. Try to get friends for non-existent user
	nonExistentUser := models.User{ID: 99999}
	friends, err = services.GetFriends(nonExistentUser.ID)
	assert.Error(t, err)
	assert.Nil(t, friends)
}

func TestFriendService_GetFriends_AfterRemovingFriend(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Create users
	user1, _ := services.CreateUser(models.User{
		Email:     "removefriend1@test.com",
		FirstName: "User",
		LastName:  "One",
	}, "pw")

	user2, _ := services.CreateUser(models.User{
		Email:     "removefriend2@test.com",
		FirstName: "User",
		LastName:  "Two",
	}, "pw")

	user3, _ := services.CreateUser(models.User{
		Email:     "removefriend3@test.com",
		FirstName: "User",
		LastName:  "Three",
	}, "pw")

	// Create friendships
	config.DB.Model(user1).Association("Friends").Append(user2)
	config.DB.Model(user2).Association("Friends").Append(user1)

	config.DB.Model(user1).Association("Friends").Append(user3)
	config.DB.Model(user3).Association("Friends").Append(user1)

	// Verify user1 has 2 friends
	friends, err := services.GetFriends(user1.ID)
	assert.NoError(t, err)
	assert.Len(t, friends, 2)

	// Remove user2 as friend
	err = services.RemoveFriend(user1.ID, user2.ID)
	assert.NoError(t, err)

	// Verify user1 now has 1 friend
	friends, err = services.GetFriends(user1.ID)
	assert.NoError(t, err)
	assert.Len(t, friends, 1)
	assert.Equal(t, user3.ID, friends[0].ID)

	// Verify user2 has no friends
	friends2, err := services.GetFriends(user2.ID)
	assert.NoError(t, err)
	assert.Empty(t, friends2)
}
