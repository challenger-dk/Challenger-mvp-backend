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
	err = services.DeleteUser(createdUser.ID)
	assert.NoError(t, err)

	// Verify Deletion
	_, err = services.GetUserByID(createdUser.ID)
	assert.Error(t, err)
}

func TestUserService_Settings(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	user, _ := services.CreateUser(models.User{Email: "settings@test.com", FirstName: "Set", LastName: "Tings", Settings: &models.UserSettings{}}, "pw")

	// 1. Get Settings
	settings, err := services.GetUserSettings(user.ID)
	assert.NoError(t, err)
	assert.True(t, settings.NotifyTeamInvite) // Default is true

	// 2. Update Settings
	f := false
	updateDto := dto.UserSettingsUpdateDto{
		NotifyTeamInvite:      &f,
		NotifyFriendReq:       &f,
		NotifyChallengeInvite: &f,
		NotifyChallengeUpdate: &f,
	}
	err = services.UpdateUserSettings(user.ID, updateDto)
	assert.NoError(t, err)

	// 3. Verify Update
	updatedSettings, _ := services.GetUserSettings(user.ID)
	assert.False(t, updatedSettings.NotifyTeamInvite)
	assert.False(t, updatedSettings.NotifyFriendReq)
}

func TestUserService_GetUsers(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	u1, _ := services.CreateUser(models.User{Email: "u1@test.com", FirstName: "U1", LastName: "L1"}, "pw")
	u2, _ := services.CreateUser(models.User{Email: "u2@test.com", FirstName: "U2", LastName: "L2"}, "pw")
	u3, _ := services.CreateUser(models.User{Email: "u3@test.com", FirstName: "U3", LastName: "L3"}, "pw")

	// 1. Get Users as u1 (Should see u2 and u3)
	users, err := services.GetUsers(u1.ID)
	assert.NoError(t, err)
	// Current impl excludes self, so should be 2
	assert.Equal(t, 2, len(users))

	// 2. Block u2 (u1 blocks u2)
	err = services.BlockUser(u1.ID, u2.ID)
	assert.NoError(t, err)

	// 3. Get Users as u1 (Should ONLY see u3)
	usersAfterBlock, err := services.GetUsers(u1.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(usersAfterBlock))
	assert.Equal(t, u3.ID, usersAfterBlock[0].ID)

	// 4. Get Users as u2 (Should ONLY see u3, u1 is hidden due to block)
	usersForBlocked, err := services.GetUsers(u2.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(usersForBlocked))
	assert.Equal(t, u3.ID, usersForBlocked[0].ID)
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
