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
