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

func TestNotifiUserMissingParticipants12H(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Freeze time
	fixed := time.Date(2026, 1, 11, 10, 0, 0, 0, time.UTC)
	oldNow := tasks.NowFunc
	tasks.NowFunc = func() time.Time { return fixed }
	defer func() { tasks.NowFunc = oldNow }()

	creator, _ := services.CreateUser(models.User{Email: "creatorMiss@test.com", FirstName: "C", LastName: "Creator"}, "pwd1")
	participant, _ := services.CreateUser(models.User{Email: "participantMiss@test.com", FirstName: "P", LastName: "Participant"}, "pwd1")
	other, _ := services.CreateUser(models.User{Email: "otherMiss@test.com", FirstName: "O", LastName: "Other"}, "pwd1")

	// Case A: participants=3, creator + participant = 2 -> missing 1 => should notify
	p := func() *int { i := 3; return &i }()
	chA := models.Challenge{
		CreatorID: creator.ID,
		Date:      fixed,
		StartTime: fixed.Add(12 * time.Hour).Add(2 * time.Minute),
		Participants: p,
	}
	createdA, err := services.CreateChallenge(chA, nil)
	assert.NoError(t, err)

	// Add one participant (creator is already auto-added)
	err = services.JoinChallenge(createdA.ID, participant.ID)
	assert.NoError(t, err)

	// Run cron
	tasks.RunNotifiUserMissingParticipantsInChallenges12H()

	// Assert notification exists for creator
	var notf models.Notification
	err = config.DB.Where("user_id = ? AND resource_id = ? AND type = ?", creator.ID, createdA.ID, models.NotifTypeChallengeMissingParticipants).First(&notf).Error
	assert.NoError(t, err)

	// Dedup: run again
	tasks.RunNotifiUserMissingParticipantsInChallenges12H()
	var cnt int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND resource_id = ? AND type = ?", creator.ID, createdA.ID, models.NotifTypeChallengeMissingParticipants).Count(&cnt)
	assert.Equal(t, int64(1), cnt)

	// Case B: participants=5, only 2 members present -> missing 3 => DO NOT notify
	p2 := func() *int { i := 5; return &i }()
	chB := models.Challenge{
		CreatorID: creator.ID,
		Date:      fixed,
		StartTime: fixed.Add(12 * time.Hour).Add(2 * time.Minute),
		Participants: p2,
	}
	createdB, err := services.CreateChallenge(chB, nil)
	assert.NoError(t, err)

	// Add one participant (creator auto-added)
	err = services.JoinChallenge(createdB.ID, other.ID)
	assert.NoError(t, err)

	// Run cron
	tasks.RunNotifiUserMissingParticipantsInChallenges12H()

	var cntB int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND resource_id = ? AND type = ?", creator.ID, createdB.ID, models.NotifTypeChallengeMissingParticipants).Count(&cntB)
	// Any missing participants (not just 1-2) should trigger notification
	assert.Equal(t, int64(1), cntB)

	// Case C: participants=nil => should be skipped
	chC := models.Challenge{
		CreatorID: creator.ID,
		Date:      fixed,
		StartTime: fixed.Add(12 * time.Hour).Add(2 * time.Minute),
		Participants: nil,
	}
	createdC, err := services.CreateChallenge(chC, nil)
	assert.NoError(t, err)

	// Add a participant so it's technically not full but there is no limit
	err = services.JoinChallenge(createdC.ID, participant.ID)
	assert.NoError(t, err)

	// Run cron
	tasks.RunNotifiUserMissingParticipantsInChallenges12H()

	var cntC int64
	config.DB.Model(&models.Notification{}).Where("user_id = ? AND resource_id = ? AND type = ?", creator.ID, createdC.ID, models.NotifTypeChallengeMissingParticipants).Count(&cntC)
	assert.Equal(t, int64(0), cntC)
}
