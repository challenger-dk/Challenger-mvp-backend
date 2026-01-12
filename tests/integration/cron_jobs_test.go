package integration

import (
	"server/api/cron/tasks"
	"server/common/config"
	"server/common/models"
	"server/common/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ------- TESTS FOR UPCOMMING CHALLENGES NOTIFICATIONS ------- \\
func TestNotifiUserUpcommingChallenges24H(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Freeze time for deterministic behavior
	fixed := time.Date(2026, 1, 11, 10, 0, 0, 0, time.UTC)
	oldNow := tasks.NowFunc
	tasks.NowFunc = func() time.Time { return fixed }
	defer func() { tasks.NowFunc = oldNow }()

	creator, _ := services.CreateUser(models.User{Email: "creator24@test.com", FirstName: "C", LastName: "Creator"}, "pwd1")
	participant, _ := services.CreateUser(models.User{Email: "participant24@test.com", FirstName: "P", LastName: "Participant"}, "pwd1")

	ch := models.Challenge{
		CreatorID: creator.ID,
		Date:      fixed,
		StartTime: fixed.Add(24 * time.Hour).Add(2 * time.Minute),
	}
	created, err := services.CreateChallenge(ch, nil)
	assert.NoError(t, err)

	// Add participant
	err = services.JoinChallenge(created.ID, participant.ID)
	assert.NoError(t, err)

	// Run the cron task (wrapper)
	tasks.RunNotifiUserUpcommingChallenges24H()

	// Assert notification exists for participant
	var n models.Notification
	err = config.DB.Where("user_id = ? AND resource_id = ? AND type = ?", participant.ID, created.ID, models.NotifTypeChallengeUpcomming24H).First(&n).Error
	assert.NoError(t, err)

	// Run again to ensure dedupe prevents duplicates
	tasks.RunNotifiUserUpcommingChallenges24H()
	var count int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND resource_id = ? AND type = ?", participant.ID, created.ID, models.NotifTypeChallengeUpcomming24H).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestNotifiUserUpcommingChallenges1H(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Freeze time for deterministic behavior
	fixed := time.Date(2026, 1, 11, 10, 0, 0, 0, time.UTC)
	oldNow := tasks.NowFunc
	tasks.NowFunc = func() time.Time { return fixed }
	defer func() { tasks.NowFunc = oldNow }()

	creator, _ := services.CreateUser(models.User{Email: "creator1@test.com", FirstName: "C", LastName: "Creator"}, "pwd1")
	participant, _ := services.CreateUser(models.User{Email: "participant1@test.com", FirstName: "P", LastName: "Participant"}, "pwd1")

	ch := models.Challenge{
		CreatorID: creator.ID,
		Date:      fixed,
		StartTime: fixed.Add(1 * time.Hour).Add(2 * time.Minute),
	}
	created, err := services.CreateChallenge(ch, nil)
	assert.NoError(t, err)

	// Add participant
	err = services.JoinChallenge(created.ID, participant.ID)
	assert.NoError(t, err)

	// Run the cron task (wrapper)
	tasks.RunNotifiUserUpcommingChallenges1H()

	// Assert notification exists for participant
	var n models.Notification
	err = config.DB.Where("user_id = ? AND resource_id = ? AND type = ?", participant.ID, created.ID, models.NotifTypeChallengeUpcomming1H).First(&n).Error
	assert.NoError(t, err)

	// Run again to ensure dedupe prevents duplicates
	tasks.RunNotifiUserUpcommingChallenges1H()
	var count int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND resource_id = ? AND type = ?", participant.ID, created.ID, models.NotifTypeChallengeUpcomming1H).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestNotifiUserInvitedToChallengeNotAnswered24H(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Freeze time for deterministic behavior
	fixed := time.Date(2026, 1, 11, 10, 0, 0, 0, time.UTC)
	oldNow := tasks.NowFunc
	tasks.NowFunc = func() time.Time { return fixed }
	defer func() { tasks.NowFunc = oldNow }()

	creator, _ := services.CreateUser(models.User{Email: "creatorInv@test.com", FirstName: "C", LastName: "Creator"}, "pwd1")
	invitee, _ := services.CreateUser(models.User{Email: "inviteeInv@test.com", FirstName: "I", LastName: "Invitee"}, "pwd1")

	ch := models.Challenge{
		CreatorID: creator.ID,
		Date:      fixed,
		StartTime: fixed.Add(24 * time.Hour).Add(2 * time.Minute),
	}
	created, err := services.CreateChallenge(ch, nil)
	assert.NoError(t, err)

	// Create a pending invitation
	inv := models.Invitation{
		InviterId:    creator.ID,
		InviteeId:    invitee.ID,
		ResourceType: models.ResourceTypeChallenge,
		ResourceID:   created.ID,
		Status:       models.StatusPending,
	}
	config.DB.Create(&inv)

	// Run the cron task (wrapper)
	tasks.RunNotifiUserInvitedToChallengeNotAnswered24H()

	// Assert notification exists for invitee
	var n models.Notification
	err = config.DB.Where("user_id = ? AND resource_id = ? AND type = ?", invitee.ID, created.ID, models.NotifTypeChallengeNotAnswered24H).First(&n).Error
	assert.NoError(t, err)

	// Run again to ensure dedupe prevents duplicates
	tasks.RunNotifiUserInvitedToChallengeNotAnswered24H()
	var count int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND resource_id = ? AND type = ?", invitee.ID, created.ID, models.NotifTypeChallengeNotAnswered24H).Count(&count)
	assert.Equal(t, int64(1), count)
}
