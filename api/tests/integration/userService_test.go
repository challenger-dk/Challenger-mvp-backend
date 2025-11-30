package integration

import (
	"server/api/services"
	"server/common/appError"
	"server/common/config"
	"server/common/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_CRUD(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// 1. Create User Success
	email := "test@example.com"
	createdUser, err := services.CreateUser(email, "password123", "John", "Doe", []string{"Tennis"})
	assert.NoError(t, err)
	assert.NotZero(t, createdUser.ID)

	// 2. Create Duplicate User (Should Fail)
	_, err = services.CreateUser(email, "password123", "John", "Doe", nil)
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

	user, _ := services.CreateUser("settings@test.com", "pw", "Set", "Tings", nil)

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

	services.CreateUser("u1@test.com", "pw", "U1", "L1", nil)
	services.CreateUser("u2@test.com", "pw", "U2", "L2", nil)

	users, err := services.GetUsers()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 2)
}

func TestUserService_FriendshipAndStats(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	u1, _ := services.CreateUser("f1@test.com", "pw", "F1", "L1", []string{"Tennis", "Football"})
	u2, _ := services.CreateUser("f2@test.com", "pw", "F2", "L2", []string{"Tennis", "Basketball"})
	u3, _ := services.CreateUser("f3@test.com", "pw", "F3", "L3", nil)

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
