package integration

import (
	"server/common/appError"
	"server/common/config"
	"server/common/dto"
	"server/common/models"
	"server/common/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_CRUD(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// 1. Create User Success
	email := "test@example.com"
	userModel := models.User{
		Email:          email,
		FirstName:      "John",
		LastName:       "Doe",
		FavoriteSports: []models.Sport{{Name: "Tennis"}},
		Settings:       &models.UserSettings{},
	}
	createdUser, err := services.CreateUser(userModel, "password123")
	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)

	// 2. Create Duplicate User (Should Fail)
	_, err = services.CreateUser(models.User{Email: email, FirstName: "John", LastName: "Doe"}, "password123")
	assert.ErrorIs(t, err, appError.ErrUserExists)

	// 3. Get User By ID (Success)
	fetchedUser, err := services.GetUserByID(createdUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, fetchedUser.ID)
	assert.Equal(t, "Tennis", fetchedUser.FavoriteSports[0].Name)

	// 4. Get User By ID (Not Found)
	_, err = services.GetUserByID(99999)
	assert.Error(t, err)

	// 5. Get User With Settings
	userWithSettings, err := services.GetUserByIDWithSettings(createdUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, userWithSettings.Settings)

	// 6. Update User Info
	updateDto := dto.UserUpdateDto{
		FirstName:      "Johnny",
		LastName:       "Doey",
		ProfilePicture: "new_pic.jpg",
		Bio:            "New Bio",
		FavoriteSports: []string{"Football"}, // Change sport
	}
	err = services.UpdateUser(createdUser.ID, updateDto)
	assert.NoError(t, err)

	// Verify Update
	updatedUser, _ := services.GetUserByID(createdUser.ID)
	assert.Equal(t, "Johnny", updatedUser.FirstName)
	assert.Equal(t, "Doey", updatedUser.LastName)
	assert.Equal(t, "new_pic.jpg", updatedUser.ProfilePicture)
	assert.Equal(t, "New Bio", updatedUser.Bio)
	assert.Len(t, updatedUser.FavoriteSports, 1)
	assert.Equal(t, "Football", updatedUser.FavoriteSports[0].Name)

	// 7. Delete User
	err = services.DeleteUser(*createdUser, email)
	assert.NoError(t, err)

	// Verify Deletion
	_, err = services.GetUserByID(createdUser.ID)
	assert.Error(t, err)

	// 8. Test Delete User with wrong email (should fail)
	user2, _ := services.CreateUser(models.User{Email: "test2@example.com", FirstName: "Jane", LastName: "Doe"}, "password123")
	err = services.DeleteUser(*user2, "wrong@email.com")
	assert.ErrorIs(t, err, appError.ErrInvalidCredentials)

	// Verify user2 still exists
	_, err = services.GetUserByID(user2.ID)
	assert.NoError(t, err)
}

func TestUserService_Settings(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	user, _ := services.CreateUser(models.User{Email: "settings@test.com", FirstName: "Set", LastName: "Tings", Settings: &models.UserSettings{}}, "pw")

	// 1. Get Settings
	settings, err := services.GetUserSettings(user.ID)
	assert.NoError(t, err)
	assert.True(t, settings.NotifyTeamInvites) // Default is true

	// 2. Update Settings
	f := false
	updateDto := dto.UserSettingsUpdateDto{
		NotifyTeamInvites:      &f,
		NotifyFriendRequests:   &f,
		NotifyChallengeInvites: &f,
		NotifyChallengeUpdates: &f,
	}
	err = services.UpdateUserSettings(user.ID, updateDto)
	assert.NoError(t, err)

	// 3. Verify Update
	updatedSettings, _ := services.GetUserSettings(user.ID)
	assert.False(t, updatedSettings.NotifyTeamInvites)
	assert.False(t, updatedSettings.NotifyFriendRequests)
}

func TestUserService_GetUsers(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	u1, _ := services.CreateUser(models.User{Email: "u1@test.com", FirstName: "Alice", LastName: "Anderson"}, "pw")
	u2, _ := services.CreateUser(models.User{Email: "u2@test.com", FirstName: "Bob", LastName: "Brown"}, "pw")
	u3, _ := services.CreateUser(models.User{Email: "u3@test.com", FirstName: "Charlie", LastName: "Chen"}, "pw")

	// 1. Get Users as u1 (Should see u2 and u3)
	users, nextCursor, err := services.GetUsers(u1.ID, "", 20, nil)
	assert.NoError(t, err)
	// Current impl excludes self, so should be 2
	assert.Equal(t, 2, len(users))
	assert.Nil(t, nextCursor) // No more pages

	// 2. Block u2 (u1 blocks u2)
	err = services.BlockUser(u1.ID, u2.ID)
	assert.NoError(t, err)

	// 3. Get Users as u1 (Should ONLY see u3)
	usersAfterBlock, nextCursor, err := services.GetUsers(u1.ID, "", 20, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(usersAfterBlock))
	assert.Equal(t, u3.ID, usersAfterBlock[0].ID)
	assert.Nil(t, nextCursor)

	// 4. Get Users as u2 (Should ONLY see u3, u1 is hidden due to block)
	usersForBlocked, nextCursor, err := services.GetUsers(u2.ID, "", 20, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(usersForBlocked))
	assert.Equal(t, u3.ID, usersForBlocked[0].ID)
	assert.Nil(t, nextCursor)
}

func TestUserService_GetUsers_Search(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Create users with searchable names
	u1, _ := services.CreateUser(models.User{Email: "alice@test.com", FirstName: "Alice", LastName: "Anderson"}, "pw")
	u2, _ := services.CreateUser(models.User{Email: "bob@test.com", FirstName: "Bob", LastName: "Brown"}, "pw")
	u3, _ := services.CreateUser(models.User{Email: "charlie@test.com", FirstName: "Charlie", LastName: "Chen"}, "pw")
	_, _ = services.CreateUser(models.User{Email: "david@test.com", FirstName: "David", LastName: "Smith"}, "pw")

	// 1. Search by first name (case-insensitive)
	users, nextCursor, err := services.GetUsers(u1.ID, "bob", 20, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, u2.ID, users[0].ID)
	assert.Nil(t, nextCursor)

	// 2. Search by last name
	users, nextCursor, err = services.GetUsers(u1.ID, "chen", 20, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, u3.ID, users[0].ID)
	assert.Nil(t, nextCursor)

	// 3. Search by partial match (case-insensitive)
	users, nextCursor, err = services.GetUsers(u1.ID, "CHAR", 20, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, u3.ID, users[0].ID) // Should find Charlie
	assert.Nil(t, nextCursor)

	// 4. Search by full name
	users, nextCursor, err = services.GetUsers(u1.ID, "charlie chen", 20, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, u3.ID, users[0].ID)
	assert.Nil(t, nextCursor)

	// 5. Search with no results
	users, nextCursor, err = services.GetUsers(u1.ID, "xyz123", 20, nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(users))
	assert.Nil(t, nextCursor)

	// 6. Empty search query (should return all users except self)
	users, nextCursor, err = services.GetUsers(u1.ID, "", 20, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(users)) // u2, u3, u4 (not u1)
	assert.Nil(t, nextCursor)

	// 7. Verify self is always excluded (search for different user to avoid this issue)
	users, nextCursor, err = services.GetUsers(u2.ID, "alice", 20, nil)
	assert.NoError(t, err)
	// Searching as Bob (u2) for Alice (u1), should find Alice
	assert.Equal(t, 1, len(users))
	assert.Equal(t, u1.ID, users[0].ID)
	assert.Nil(t, nextCursor)

	// 8. Verify requester is excluded when searching for themselves
	users, nextCursor, err = services.GetUsers(u1.ID, "anderson", 20, nil)
	assert.NoError(t, err)
	// Alice Anderson (u1) is the requester, should not appear
	for _, user := range users {
		assert.NotEqual(t, u1.ID, user.ID, "Requester should never appear in search results")
	}
	assert.Nil(t, nextCursor)
}

func TestUserService_GetUsers_Pagination(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Create requesting user
	requester, _ := services.CreateUser(models.User{Email: "requester@test.com", FirstName: "Requester", LastName: "User"}, "pw")

	// Create 5 users (ordered by last_name, first_name, id)
	u1, _ := services.CreateUser(models.User{Email: "u1@test.com", FirstName: "Alice", LastName: "Anderson"}, "pw")
	u2, _ := services.CreateUser(models.User{Email: "u2@test.com", FirstName: "Bob", LastName: "Brown"}, "pw")
	u3, _ := services.CreateUser(models.User{Email: "u3@test.com", FirstName: "Charlie", LastName: "Chen"}, "pw")
	u4, _ := services.CreateUser(models.User{Email: "u4@test.com", FirstName: "Diana", LastName: "Davis"}, "pw")
	u5, _ := services.CreateUser(models.User{Email: "u5@test.com", FirstName: "Eve", LastName: "Evans"}, "pw")

	// 1. Get first page (limit 2)
	page1, nextCursor, err := services.GetUsers(requester.ID, "", 2, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(page1))
	assert.NotNil(t, nextCursor)        // Should have next page
	assert.Equal(t, u1.ID, page1[0].ID) // Alice Anderson
	assert.Equal(t, u2.ID, page1[1].ID) // Bob Brown

	// 2. Get second page using cursor
	page2, nextCursor, err := services.GetUsers(requester.ID, "", 2, nextCursor)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(page2))
	assert.NotNil(t, nextCursor)        // Should have next page
	assert.Equal(t, u3.ID, page2[0].ID) // Charlie Chen
	assert.Equal(t, u4.ID, page2[1].ID) // Diana Davis

	// 3. Get third page (last page)
	page3, nextCursor, err := services.GetUsers(requester.ID, "", 2, nextCursor)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(page3))
	assert.Nil(t, nextCursor)           // No more pages
	assert.Equal(t, u5.ID, page3[0].ID) // Eve Evans

	// 4. Test pagination with search
	page1, nextCursor, err = services.GetUsers(requester.ID, "e", 2, nil)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(page1), 1) // Should find users with 'e' in name

	// 5. Test limit clamping (max 50)
	users, nextCursor, err := services.GetUsers(requester.ID, "", 100, nil)
	assert.NoError(t, err)
	assert.Equal(t, 5, len(users)) // All 5 users
	assert.Nil(t, nextCursor)
}

func TestUserService_Blocking(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	u1, _ := services.CreateUser(models.User{Email: "b1@test.com", FirstName: "B1", LastName: "L1"}, "pw")
	u2, _ := services.CreateUser(models.User{Email: "b2@test.com", FirstName: "B2", LastName: "L2"}, "pw")

	// 1. Check Visibility (Should be visible)
	user, err := services.GetVisibleUser(u1.ID, u2.ID)
	assert.NoError(t, err)
	assert.Equal(t, u2.ID, user.ID)

	// 2. Block User
	err = services.BlockUser(u1.ID, u2.ID)
	assert.NoError(t, err)

	// 3. Check Visibility (Should be Not Found)
	_, err = services.GetVisibleUser(u1.ID, u2.ID)
	assert.ErrorIs(t, err, appError.ErrUserNotFound)

	// 4. Check Stats (Should be Not Found)
	_, err = services.GetInCommonStats(u1.ID, u2.ID)
	assert.ErrorIs(t, err, appError.ErrUserNotFound)

	// 5. Unblock
	err = services.UnblockUser(u1.ID, u2.ID)
	assert.NoError(t, err)

	// 6. Check Visibility (Should be visible again)
	user, err = services.GetVisibleUser(u1.ID, u2.ID)
	assert.NoError(t, err)
	assert.Equal(t, u2.ID, user.ID)
}

func TestUserService_FriendshipAndStats(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	u1, _ := services.CreateUser(models.User{
		Email:          "f1@test.com",
		FirstName:      "F1",
		LastName:       "L1",
		FavoriteSports: []models.Sport{{Name: "Tennis"}, {Name: "Football"}},
	}, "pw")

	u2, _ := services.CreateUser(models.User{
		Email:          "f2@test.com",
		FirstName:      "F2",
		LastName:       "L2",
		FavoriteSports: []models.Sport{{Name: "Tennis"}, {Name: "Basketball"}},
	}, "pw")

	u3, _ := services.CreateUser(models.User{Email: "f3@test.com", FirstName: "F3", LastName: "L3"}, "pw")

	// 1. Manually create friendship between u1 and u2 (Direct Friendship)
	config.DB.Model(u1).Association("Friends").Append(u2)
	config.DB.Model(u2).Association("Friends").Append(u1)

	// 2. Test GetInCommonStats
	stats, err := services.GetInCommonStats(u1.ID, u2.ID)
	assert.NoError(t, err)

	// IMPORTANT: Mutual friends count should be 0 because they don't have a 3rd friend in common yet.
	assert.Equal(t, int64(0), stats.CommonFriendsCount)

	// They DO have 1 common sport (Tennis)
	assert.Len(t, stats.CommonSports, 1)
	assert.Equal(t, "Tennis", stats.CommonSports[0].Name)

	// 3. Add u3 as a mutual friend to BOTH
	config.DB.Model(u1).Association("Friends").Append(u3)
	config.DB.Model(u2).Association("Friends").Append(u3)

	stats, _ = services.GetInCommonStats(u1.ID, u2.ID)

	// NOW they have 1 mutual friend (u3)
	assert.Equal(t, int64(1), stats.CommonFriendsCount)

	// 4. Remove Friend (u1 removes u2)
	err = services.RemoveFriend(u1.ID, u2.ID)
	assert.NoError(t, err)

	// Verify removal in DB
	var friendCount int64
	config.DB.Table("user_friends").Where("user_id = ? AND friend_id = ?", u1.ID, u2.ID).Count(&friendCount)
	assert.Equal(t, int64(0), friendCount)

	// 5. Remove Friend Error (Same ID)
	err = services.RemoveFriend(u1.ID, u1.ID)
	assert.ErrorIs(t, err, appError.ErrInvalidFriendship)
}
