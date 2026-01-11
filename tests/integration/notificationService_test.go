package integration

import (
	"server/common/config"
	"server/common/models"
	"server/common/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotificationService_GetMyNotifications(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Create users
	user1, _ := services.CreateUser(models.User{
		Email:     "notif1@test.com",
		FirstName: "Notif",
		LastName:  "One",
	}, "pw")

	user2, _ := services.CreateUser(models.User{
		Email:     "notif2@test.com",
		FirstName: "Notif",
		LastName:  "Two",
	}, "pw")

	// Create notifications for user1
	notif1 := models.Notification{
		UserID:     user1.ID,
		Type:       models.NotifTypeFriendReq,
		Title:      "Friend Request 1",
		Content:    "You have a friend request",
		IsRead:     false,
		IsRelevant: true,
	}
	config.DB.Create(&notif1)

	notif2 := models.Notification{
		UserID:     user1.ID,
		Type:       models.NotifTypeTeamInvite,
		Title:      "Team Invite 1",
		Content:    "You are invited to a team",
		IsRead:     true,
		IsRelevant: true,
	}
	config.DB.Create(&notif2)

	notif3 := models.Notification{
		UserID:     user1.ID,
		Type:       models.NotifTypeChallengeReq,
		Title:      "Challenge Request 1",
		Content:    "You are invited to a challenge",
		IsRead:     false,
		IsRelevant: true,
	}
	config.DB.Create(&notif3)

	// Create irrelevant notification (should not appear)
	notif4 := models.Notification{
		UserID:     user1.ID,
		Type:       models.NotifTypeFriendReq,
		Title:      "Hidden Notification",
		Content:    "This should not appear",
		IsRead:     false,
		IsRelevant: false,
	}
	config.DB.Create(&notif4)

	// Create notification for user2 (should not appear)
	notif5 := models.Notification{
		UserID:     user2.ID,
		Type:       models.NotifTypeFriendReq,
		Title:      "Other User Notification",
		Content:    "This is for user2",
		IsRead:     false,
		IsRelevant: true,
	}
	config.DB.Create(&notif5)

	// 1. Get all notifications for user1 (no filters)
	allNotifs, err := services.GetMyNotifications(user1.ID, services.NotificationFilters{})
	assert.NoError(t, err)
	// Should return only relevant notifications (excluding irrelevant and user2's)
	assert.GreaterOrEqual(t, len(allNotifs), 3, "Should return at least 3 relevant notifications")

	// Verify all returned notifications are relevant and belong to user1
	for _, notif := range allNotifs {
		assert.True(t, notif.IsRelevant, "All returned notifications should be relevant")
		assert.Equal(t, user1.ID, notif.UserID, "All notifications should belong to user1")
	}

	// Verify notifications are ordered by created_at desc (most recent first)
	if len(allNotifs) >= 2 {
		for i := 0; i < len(allNotifs)-1; i++ {
			assert.True(t, allNotifs[i].CreatedAt.After(allNotifs[i+1].CreatedAt) || allNotifs[i].CreatedAt.Equal(allNotifs[i+1].CreatedAt),
				"Notifications should be ordered by created_at desc")
		}
	}

	// Verify all returned notifications are relevant and belong to user1
	for _, notif := range allNotifs {
		assert.True(t, notif.IsRelevant, "All returned notifications should be relevant")
		assert.Equal(t, user1.ID, notif.UserID, "All notifications should belong to user1")
	}

	// Explicitly mark notif4 as irrelevant and verify it's excluded
	config.DB.Model(&notif4).Update("is_relevant", false)

	// Get notifications again and verify notif4 is excluded
	allNotifsAfterUpdate, err := services.GetMyNotifications(user1.ID, services.NotificationFilters{})
	assert.NoError(t, err)
	for _, notif := range allNotifsAfterUpdate {
		assert.True(t, notif.IsRelevant, "All returned notifications should be relevant")
		if notif.ID == notif4.ID {
			t.Errorf("Irrelevant notification (ID: %d) should not be included", notif4.ID)
		}
	}

	// 2. Get only unread notifications (isRead = false means unread)
	unreadFilterFalse := false
	unreadNotifs, err := services.GetMyNotifications(user1.ID, services.NotificationFilters{
		IsRead: &unreadFilterFalse,
	})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(unreadNotifs), 2, "Should return at least 2 unread notifications")

	// Verify all returned notifications are unread and relevant
	for _, notif := range unreadNotifs {
		assert.False(t, notif.IsRead, "All returned notifications should be unread")
		assert.True(t, notif.IsRelevant, "All returned notifications should be relevant")
		assert.Equal(t, user1.ID, notif.UserID, "All notifications should belong to user1")
	}

	// 3. Get only read notifications
	readFilter := true
	readNotifs, err := services.GetMyNotifications(user1.ID, services.NotificationFilters{
		IsRead: &readFilter,
	})
	assert.NoError(t, err)
	assert.Len(t, readNotifs, 1, "Should return 1 read notification")
	assert.True(t, readNotifs[0].IsRead)

	// 4. Test limit
	limitedNotifs, err := services.GetMyNotifications(user1.ID, services.NotificationFilters{
		Limit: 2,
	})
	assert.NoError(t, err)
	assert.Len(t, limitedNotifs, 2, "Should return only 2 notifications when limit is 2")

	// 5. Test offset
	offsetNotifs, err := services.GetMyNotifications(user1.ID, services.NotificationFilters{
		Limit:  2,
		Offset: 1,
	})
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(offsetNotifs), 2, "Should return at most 2 notifications with offset")
	// Should skip the first notification(s)
	if len(offsetNotifs) > 0 && len(allNotifs) > 0 {
		assert.NotEqual(t, allNotifs[0].ID, offsetNotifs[0].ID, "First notification should be skipped due to offset")
	}

	// 6. Test default limit (50)
	manyNotifs, err := services.GetMyNotifications(user1.ID, services.NotificationFilters{
		Limit: 0, // Should default to 50
	})
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(manyNotifs), 50, "Should not exceed default limit of 50")

	// 7. Get notifications for user2 (should only get user2's notifications)
	user2Notifs, err := services.GetMyNotifications(user2.ID, services.NotificationFilters{})
	assert.NoError(t, err)
	// Verify all notifications belong to user2
	for _, notif := range user2Notifs {
		assert.Equal(t, user2.ID, notif.UserID, "All notifications should belong to user2")
		assert.True(t, notif.IsRelevant, "All notifications should be relevant")
	}
}

func TestNotificationService_MarkNotificationAsRead(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Create users
	user1, _ := services.CreateUser(models.User{
		Email:     "markread1@test.com",
		FirstName: "Mark",
		LastName:  "One",
	}, "pw")

	user2, _ := services.CreateUser(models.User{
		Email:     "markread2@test.com",
		FirstName: "Mark",
		LastName:  "Two",
	}, "pw")

	// Create notifications
	notif1 := models.Notification{
		UserID:     user1.ID,
		Type:       models.NotifTypeFriendReq,
		Title:      "Notification 1",
		Content:    "Test notification",
		IsRead:     false,
		IsRelevant: true,
	}
	config.DB.Create(&notif1)

	notif2 := models.Notification{
		UserID:     user1.ID,
		Type:       models.NotifTypeTeamInvite,
		Title:      "Notification 2",
		Content:    "Test notification 2",
		IsRead:     false,
		IsRelevant: true,
	}
	config.DB.Create(&notif2)

	notif3 := models.Notification{
		UserID:     user2.ID,
		Type:       models.NotifTypeFriendReq,
		Title:      "User2 Notification",
		Content:    "This is for user2",
		IsRead:     false,
		IsRelevant: true,
	}
	config.DB.Create(&notif3)

	// 1. Mark notification as read successfully
	err := services.MarkNotificationAsRead(notif1.ID, user1.ID)
	assert.NoError(t, err)

	// Verify notification is marked as read
	var updatedNotif1 models.Notification
	config.DB.First(&updatedNotif1, notif1.ID)
	assert.True(t, updatedNotif1.IsRead, "Notification should be marked as read")

	// Verify other notifications are still unread
	var updatedNotif2 models.Notification
	config.DB.First(&updatedNotif2, notif2.ID)
	assert.False(t, updatedNotif2.IsRead, "Other notification should remain unread")

	// 2. Try to mark another user's notification (should fail silently or succeed but not affect)
	err = services.MarkNotificationAsRead(notif3.ID, user1.ID)
	// This should not error, but also should not mark user2's notification as read
	var updatedNotif3 models.Notification
	config.DB.First(&updatedNotif3, notif3.ID)
	// The update should not affect user2's notification since the WHERE clause includes user_id
	assert.False(t, updatedNotif3.IsRead, "User2's notification should not be affected")

	// 3. Mark notification as read again (idempotent)
	err = services.MarkNotificationAsRead(notif1.ID, user1.ID)
	assert.NoError(t, err)
	config.DB.First(&updatedNotif1, notif1.ID)
	assert.True(t, updatedNotif1.IsRead, "Notification should remain read")

	// 4. Mark second notification as read
	err = services.MarkNotificationAsRead(notif2.ID, user1.ID)
	assert.NoError(t, err)
	config.DB.First(&updatedNotif2, notif2.ID)
	assert.True(t, updatedNotif2.IsRead, "Second notification should be marked as read")

	// 5. Try to mark non-existent notification
	err = services.MarkNotificationAsRead(99999, user1.ID)
	assert.NoError(t, err) // Should not error, just updates 0 rows
}

func TestNotificationService_MarkAllNotificationsAsRead(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Create users
	user1, _ := services.CreateUser(models.User{
		Email:     "markall1@test.com",
		FirstName: "MarkAll",
		LastName:  "One",
	}, "pw")

	user2, _ := services.CreateUser(models.User{
		Email:     "markall2@test.com",
		FirstName: "MarkAll",
		LastName:  "Two",
	}, "pw")

	// Create multiple unread notifications for user1
	for i := 0; i < 5; i++ {
		notif := models.Notification{
			UserID:     user1.ID,
			Type:       models.NotifTypeFriendReq,
			Title:      "Notification " + string(rune(i)),
			Content:    "Test notification",
			IsRead:     false,
			IsRelevant: true,
		}
		config.DB.Create(&notif)
	}

	// Create one read notification for user1
	readNotif := models.Notification{
		UserID:     user1.ID,
		Type:       models.NotifTypeTeamInvite,
		Title:      "Already Read",
		Content:    "This is already read",
		IsRead:     true,
		IsRelevant: true,
	}
	config.DB.Create(&readNotif)

	// Create unread notifications for user2
	for i := 0; i < 3; i++ {
		notif := models.Notification{
			UserID:     user2.ID,
			Type:       models.NotifTypeFriendReq,
			Title:      "User2 Notification " + string(rune(i)),
			Content:    "Test notification",
			IsRead:     false,
			IsRelevant: true,
		}
		config.DB.Create(&notif)
	}

	// 1. Mark all notifications as read for user1
	err := services.MarkAllNotificationsAsRead(user1.ID)
	assert.NoError(t, err)

	// Verify all unread notifications for user1 are now read
	var user1Notifs []models.Notification
	config.DB.Where("user_id = ? AND is_relevant = ?", user1.ID, true).Find(&user1Notifs)
	for _, notif := range user1Notifs {
		assert.True(t, notif.IsRead, "All user1 notifications should be marked as read")
	}

	// Verify user2's notifications are still unread
	var user2Notifs []models.Notification
	config.DB.Where("user_id = ? AND is_relevant = ?", user2.ID, true).Find(&user2Notifs)
	for _, notif := range user2Notifs {
		assert.False(t, notif.IsRead, "User2's notifications should remain unread")
	}

	// 2. Mark all notifications as read again (idempotent)
	err = services.MarkAllNotificationsAsRead(user1.ID)
	assert.NoError(t, err)
	// Should not cause any issues

	// 3. Mark all notifications for user2
	err = services.MarkAllNotificationsAsRead(user2.ID)
	assert.NoError(t, err)

	// Verify user2's notifications are now read
	config.DB.Where("user_id = ? AND is_relevant = ?", user2.ID, true).Find(&user2Notifs)
	for _, notif := range user2Notifs {
		assert.True(t, notif.IsRead, "All user2 notifications should be marked as read")
	}

	// 4. Mark all for user with no unread notifications
	user3, _ := services.CreateUser(models.User{
		Email:     "markall3@test.com",
		FirstName: "MarkAll",
		LastName:  "Three",
	}, "pw")
	err = services.MarkAllNotificationsAsRead(user3.ID)
	assert.NoError(t, err) // Should not error even if no notifications exist
}

func TestNotificationService_CreateTeamDeletedNotification(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Create users
	creator, _ := services.CreateUser(models.User{
		Email:     "teamdel1@test.com",
		FirstName: "TeamDel",
		LastName:  "Creator",
	}, "pw")

	member1, _ := services.CreateUser(models.User{
		Email:     "teamdel2@test.com",
		FirstName: "TeamDel",
		LastName:  "Member1",
	}, "pw")

	member2, _ := services.CreateUser(models.User{
		Email:     "teamdel3@test.com",
		FirstName: "TeamDel",
		LastName:  "Member2",
	}, "pw")

	// Create a team
	team, _ := services.CreateTeam(models.Team{
		Name:      "Team To Delete",
		CreatorID: creator.ID,
	})

	// Add members to team
	var teamModel models.Team
	config.DB.First(&teamModel, team.ID)
	config.DB.Model(&teamModel).Association("Users").Append(member1)
	config.DB.Model(&teamModel).Association("Users").Append(member2)

	// 1. Create team deleted notification for member1
	services.CreateTeamDeletedNotification(config.DB, *member1, teamModel)

	// Verify notification was created
	var notif1 models.Notification
	err := config.DB.Where("user_id = ? AND type = ?", member1.ID, models.NotifTypeTeamDeleted).First(&notif1).Error
	assert.NoError(t, err)
	assert.Equal(t, "Et hold du er medlem af er blevet slettet", notif1.Title)
	assert.Contains(t, notif1.Content, "Team To Delete")
	assert.Equal(t, models.NotifTypeTeamDeleted, notif1.Type)
	assert.False(t, notif1.IsRead)
	assert.True(t, notif1.IsRelevant)

	// 2. Create team deleted notification for member2
	services.CreateTeamDeletedNotification(config.DB, *member2, teamModel)

	// Verify notification was created
	var notif2 models.Notification
	err = config.DB.Where("user_id = ? AND type = ?", member2.ID, models.NotifTypeTeamDeleted).First(&notif2).Error
	assert.NoError(t, err)
	assert.Equal(t, member2.ID, notif2.UserID)
	assert.Equal(t, models.NotifTypeTeamDeleted, notif2.Type)

	// 3. Verify creator did not get notification (only members do)
	var creatorNotifs []models.Notification
	config.DB.Where("user_id = ? AND type = ?", creator.ID, models.NotifTypeTeamDeleted).Find(&creatorNotifs)
	assert.Empty(t, creatorNotifs, "Creator should not receive team deleted notification")

	// 4. Create notification for multiple team members
	team2, _ := services.CreateTeam(models.Team{
		Name:      "Another Team",
		CreatorID: creator.ID,
	})

	var teamModel2 models.Team
	config.DB.First(&teamModel2, team2.ID)
	config.DB.Model(&teamModel2).Association("Users").Append(member1)
	config.DB.Model(&teamModel2).Association("Users").Append(member2)

	// Create notifications for all members
	services.CreateTeamDeletedNotification(config.DB, *member1, teamModel2)
	services.CreateTeamDeletedNotification(config.DB, *member2, teamModel2)

	// Verify both members got notifications
	var member1Notifs []models.Notification
	config.DB.Where("user_id = ? AND type = ?", member1.ID, models.NotifTypeTeamDeleted).Find(&member1Notifs)
	assert.GreaterOrEqual(t, len(member1Notifs), 2, "Member1 should have at least 2 team deleted notifications")

	var member2Notifs []models.Notification
	config.DB.Where("user_id = ? AND type = ?", member2.ID, models.NotifTypeTeamDeleted).Find(&member2Notifs)
	assert.GreaterOrEqual(t, len(member2Notifs), 2, "Member2 should have at least 2 team deleted notifications")

	// 5. Test with different team name
	team3, _ := services.CreateTeam(models.Team{
		Name:      "Special Team Name",
		CreatorID: creator.ID,
	})
	var teamModel3 models.Team
	config.DB.First(&teamModel3, team3.ID)
	config.DB.Model(&teamModel3).Association("Users").Append(member1)

	services.CreateTeamDeletedNotification(config.DB, *member1, teamModel3)

	var notif3 models.Notification
	config.DB.Where("user_id = ? AND type = ? AND content LIKE ?", member1.ID, models.NotifTypeTeamDeleted, "%Special Team Name%").First(&notif3)
	assert.NotZero(t, notif3.ID)
	assert.Contains(t, notif3.Content, "Special Team Name")
}

func TestNotificationService_GetMyNotifications_WithActorPreload(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Create users
	actor, _ := services.CreateUser(models.User{
		Email:     "actor@test.com",
		FirstName: "Actor",
		LastName:  "User",
	}, "pw")

	recipient, _ := services.CreateUser(models.User{
		Email:     "recipient@test.com",
		FirstName: "Recipient",
		LastName:  "User",
	}, "pw")

	// Create notification with actor
	notif := models.Notification{
		UserID:     recipient.ID,
		ActorID:    &actor.ID,
		Type:       models.NotifTypeFriendReq,
		Title:      "Friend Request",
		Content:    "You have a friend request",
		IsRead:     false,
		IsRelevant: true,
	}
	config.DB.Create(&notif)

	// Get notifications - should preload actor
	notifs, err := services.GetMyNotifications(recipient.ID, services.NotificationFilters{})
	assert.NoError(t, err)
	assert.Len(t, notifs, 1)
	assert.NotNil(t, notifs[0].Actor, "Actor should be preloaded")
	assert.Equal(t, actor.ID, notifs[0].Actor.ID)
	assert.Equal(t, "actor@test.com", notifs[0].Actor.Email)
}
