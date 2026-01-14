package integration

import (
	"server/common/appError"
	"server/common/config"
	"server/common/dto"
	"server/common/models"
	"server/common/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChallengeService_CRUD(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	creatorModel := models.User{
		Email:     "chal@test.com",
		FirstName: "C",
		LastName:  "C",
	}
	creator, _ := services.CreateUser(creatorModel, "pw")

	// 1. Create
	chalDto := dto.ChallengeCreateDto{
		Name:      "Match",
		Sport:     "Tennis",
		Date:      time.Now(),
		StartTime: time.Now(),
		Location: dto.LocationCreateDto{
			Address: "A", Latitude: 1, Longitude: 1, PostalCode: "1", City: "C", Country: "D",
		},
	}
	model := dto.ChallengeCreateDtoToModel(chalDto)
	model.CreatorID = creator.ID

	created, err := services.CreateChallenge(model, []uint{})
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	// 2. Get All
	list, err := services.GetChallenges()
	assert.NoError(t, err)
	assert.NotEmpty(t, list)

	// 3. Get By ID
	fetched, err := services.GetChallengeByID(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, fetched.ID)

	// 4. Update
	updateModel := models.Challenge{
		Name:        "Updated Match",
		Description: "New Desc",
		Sport:       "Football",
		TeamSize:    func() *int { i := 5; return &i }(),
	}
	err = services.UpdateChallenge(created.ID, updateModel)
	assert.NoError(t, err)

	updated, _ := services.GetChallengeByID(created.ID)
	assert.Equal(t, "Updated Match", updated.Name)
	assert.Equal(t, "New Desc", updated.Description)
	assert.Equal(t, "Football", updated.Sport)
	assert.Equal(t, 5, *updated.TeamSize)

	// 5. Delete
	err = services.DeleteChallenge(created.ID)
	assert.NoError(t, err)

	_, err = services.GetChallengeByID(created.ID)
	assert.Error(t, err)
}

func TestChallengeService_Participation(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	c1, _ := services.CreateUser(models.User{Email: "c1@c.com", FirstName: "C", LastName: "C"}, "pw")
	u1, _ := services.CreateUser(models.User{Email: "u1@c.com", FirstName: "U", LastName: "U"}, "pw")

	chal := models.Challenge{
		Name:      "Join Test",
		CreatorID: c1.ID,
		Date:      time.Now(),
		StartTime: time.Now(),
		Location:  models.Location{Address: "L", Coordinates: models.Point{Lat: 0, Lon: 0}, PostalCode: "1", City: "C", Country: "C"},
	}
	created, _ := services.CreateChallenge(chal, nil)

	// 1. Join
	err := services.JoinChallenge(created.ID, u1.ID)
	assert.NoError(t, err)

	fetched, _ := services.GetChallengeByID(created.ID)
	assert.Len(t, fetched.Users, 2) // Creator + u1

	// 2. Leave
	err = services.LeaveChallenge(created.ID, u1.ID)
	assert.NoError(t, err)

	fetched, _ = services.GetChallengeByID(created.ID)
	assert.Len(t, fetched.Users, 1)
}

func TestChallengeService_FullParticipation(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	creator, _ := services.CreateUser(models.User{Email: "creatorFull@test.com", FirstName: "C", LastName: "C"}, "pw")
	user1, _ := services.CreateUser(models.User{Email: "user1full@test.com", FirstName: "U1", LastName: "One"}, "pw")
	user2, _ := services.CreateUser(models.User{Email: "user2full@test.com", FirstName: "U2", LastName: "Two"}, "pw")

	t.Logf("created users: creator=%d user1=%d user2=%d", creator.ID, user1.ID, user2.ID)

	participants := func() *int { i := 2; return &i }()

	chal := models.Challenge{
		Name:         "Full Test",
		CreatorID:    creator.ID,
		Date:         time.Now(),
		StartTime:    time.Now(),
		Location:     models.Location{Address: "L", Coordinates: models.Point{Lat: 0, Lon: 0}, PostalCode: "1", City: "C", Country: "C"},
		Participants: participants,
	}

	created, _ := services.CreateChallenge(chal, nil)

	fetchedBeforeJoin, _ := services.GetChallengeByID(created.ID)
	idsBefore := make([]uint, len(fetchedBeforeJoin.Users))
	for i, u := range fetchedBeforeJoin.Users {
		idsBefore[i] = u.ID
	}
	t.Logf("users before join: %v", idsBefore)

	// user1 joins - should fill challenge (creator + user1 = 2)
	err := services.JoinChallenge(created.ID, user1.ID)
	assert.NoError(t, err)

	fetched, _ := services.GetChallengeByID(created.ID)
	assert.Len(t, fetched.Users, 2)
	assert.Equal(t, models.ChallengeConfirmed, fetched.Status)

	// creator should get full participation notification
	var notif models.Notification
	err = config.DB.Where("user_id = ? AND type = ?", creator.ID, models.NotifTypeChallengeFullParticipation).First(&notif).Error
	assert.NoError(t, err)

	// user1 should have a joined notification
	var notif2 models.Notification
	err = config.DB.Where("user_id = ? AND type = ?", user1.ID, models.NotifTypeChallengeJoin).First(&notif2).Error
	assert.NoError(t, err)

	// user2 trying to join should error
	err = services.JoinChallenge(created.ID, user2.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, appError.ErrChallengeFullParticipation.Error())
}

func TestChallengeService_CreateWithComplexInvites(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	creator, _ := services.CreateUser(models.User{Email: "creator@c.com", FirstName: "C", LastName: "C"}, "pw")

	// User 1: Has no prior invite (Normal flow)
	u1, _ := services.CreateUser(models.User{Email: "u1@c.com", FirstName: "U1", LastName: "U1"}, "pw")

	// User 2: Has a declined invite (Should resend)
	// Note: In CreateChallenge for a NEW challenge, prior invites to *this specific challenge ID*
	// are impossible unless IDs are reused. This test primarily verifies the creation of invites
	// for a list of users.
	u2, _ := services.CreateUser(models.User{Email: "u2@c.com", FirstName: "U2", LastName: "U2"}, "pw")

	chalModel := models.Challenge{
		Name:      "Invite Match",
		CreatorID: creator.ID,
		Date:      time.Now(),
		StartTime: time.Now(),
		Location:  models.Location{Address: "L", Coordinates: models.Point{Lat: 0, Lon: 0}, PostalCode: "1", City: "C", Country: "C"},
	}

	// Invite u1 and u2
	_, err := services.CreateChallenge(chalModel, []uint{u1.ID, u2.ID, creator.ID}) // Including creator to ensure skip logic
	assert.NoError(t, err)

	// Verify invites
	inv1, _ := services.GetInvitationsByUserId(u1.ID)
	assert.Len(t, inv1, 1)
	assert.Equal(t, models.ResourceTypeChallenge, inv1[0].ResourceType)

	inv2, _ := services.GetInvitationsByUserId(u2.ID)
	assert.Len(t, inv2, 1)

	// Verify creator (should not have invite)
	invCreator, _ := services.GetInvitationsByUserId(creator.ID)
	assert.Len(t, invCreator, 0)
}

func TestChallengeService_UpdateChallengeStatusIfExpired(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	creator, _ := services.CreateUser(models.User{Email: "expire@test.com", FirstName: "Expire", LastName: "User"}, "pw")

	// 1. Create challenge with EndTime in the past (should be marked as completed)
	pastTime := time.Now().Add(-1 * time.Hour)
	chalPast := models.Challenge{
		Name:      "Expired Challenge",
		CreatorID: creator.ID,
		Date:      time.Now(),
		StartTime: time.Now(),
		EndTime:   &pastTime,
		Status:    models.ChallengeStatusOpen,
		Location:  models.Location{Address: "L", Coordinates: models.Point{Lat: 0, Lon: 0}, PostalCode: "1", City: "C", Country: "C"},
	}
	createdPast, _ := services.CreateChallenge(chalPast, nil)

	// Get challenge - should automatically update status to completed
	fetchedPast, err := services.GetChallengeByID(createdPast.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.ChallengeStatusCompleted, fetchedPast.Status, "Challenge with past EndTime should be marked as completed")

	// 2. Create challenge with EndTime in the future (should remain open)
	futureTime := time.Now().Add(1 * time.Hour)
	chalFuture := models.Challenge{
		Name:      "Future Challenge",
		CreatorID: creator.ID,
		Date:      time.Now(),
		StartTime: time.Now(),
		EndTime:   &futureTime,
		Status:    models.ChallengeStatusOpen,
		Location:  models.Location{Address: "L", Coordinates: models.Point{Lat: 0, Lon: 0}, PostalCode: "1", City: "C", Country: "C"},
	}
	createdFuture, _ := services.CreateChallenge(chalFuture, nil)

	// Get challenge - should remain open
	fetchedFuture, err := services.GetChallengeByID(createdFuture.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.ChallengeStatusOpen, fetchedFuture.Status, "Challenge with future EndTime should remain open")

	// 3. Create challenge with no EndTime (should not change status)
	chalNoEnd := models.Challenge{
		Name:      "No End Challenge",
		CreatorID: creator.ID,
		Date:      time.Now(),
		StartTime: time.Now(),
		EndTime:   nil,
		Status:    models.ChallengeStatusOpen,
		Location:  models.Location{Address: "L", Coordinates: models.Point{Lat: 0, Lon: 0}, PostalCode: "1", City: "C", Country: "C"},
	}
	createdNoEnd, _ := services.CreateChallenge(chalNoEnd, nil)

	// Get challenge - should remain open
	fetchedNoEnd, err := services.GetChallengeByID(createdNoEnd.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.ChallengeStatusOpen, fetchedNoEnd.Status, "Challenge with no EndTime should remain unchanged")

	// 4. Create challenge already marked as completed (should not change)
	pastTime2 := time.Now().Add(-2 * time.Hour)
	chalCompleted := models.Challenge{
		Name:      "Already Completed",
		CreatorID: creator.ID,
		Date:      time.Now(),
		StartTime: time.Now(),
		EndTime:   &pastTime2,
		Status:    models.ChallengeStatusCompleted,
		Location:  models.Location{Address: "L", Coordinates: models.Point{Lat: 0, Lon: 0}, PostalCode: "1", City: "C", Country: "C"},
	}
	createdCompleted, _ := services.CreateChallenge(chalCompleted, nil)

	// Get challenge - should remain completed
	fetchedCompleted, err := services.GetChallengeByID(createdCompleted.ID)
	assert.NoError(t, err)
	assert.Equal(t, models.ChallengeStatusCompleted, fetchedCompleted.Status, "Already completed challenge should remain completed")

	// 5. Test GetChallenges also updates expired challenges
	pastTime3 := time.Now().Add(-30 * time.Minute)
	chalExpired := models.Challenge{
		Name:      "Expired in List",
		CreatorID: creator.ID,
		Date:      time.Now(),
		StartTime: time.Now(),
		EndTime:   &pastTime3,
		Status:    models.ChallengeStatusReady,
		Location:  models.Location{Address: "L", Coordinates: models.Point{Lat: 0, Lon: 0}, PostalCode: "1", City: "C", Country: "C"},
	}
	createdExpired, _ := services.CreateChallenge(chalExpired, nil)

	// Get all challenges - expired one should be updated
	allChallenges, err := services.GetChallenges()
	assert.NoError(t, err)

	var foundExpired *models.Challenge
	for i := range allChallenges {
		if allChallenges[i].ID == createdExpired.ID {
			foundExpired = &allChallenges[i]
			break
		}
	}
	assert.NotNil(t, foundExpired, "Expired challenge should be in the list")
	assert.Equal(t, models.ChallengeStatusCompleted, foundExpired.Status, "Expired challenge should be marked as completed in GetChallenges")
}

func TestChallengeService_AddUserToChallenge(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	creator, _ := services.CreateUser(models.User{Email: "adduser@test.com", FirstName: "Add", LastName: "User"}, "pw")
	user1, _ := services.CreateUser(models.User{Email: "user1add@test.com", FirstName: "User", LastName: "One"}, "pw")
	user2, _ := services.CreateUser(models.User{Email: "user2add@test.com", FirstName: "User", LastName: "Two"}, "pw")

	// Create a challenge
	chal := models.Challenge{
		Name:      "Add User Test",
		CreatorID: creator.ID,
		Date:      time.Now(),
		StartTime: time.Now(),
		Location:  models.Location{Address: "L", Coordinates: models.Point{Lat: 0, Lon: 0}, PostalCode: "1", City: "C", Country: "C"},
	}
	created, _ := services.CreateChallenge(chal, nil)

	// Verify creator is already in challenge
	fetched, _ := services.GetChallengeByID(created.ID)
	assert.Len(t, fetched.Users, 1, "Creator should be in challenge")
	assert.Equal(t, creator.ID, fetched.Users[0].ID)

	// Test addUserToChallenge indirectly through invitation acceptance
	// 1. Create invitation for user1
	invitation := models.Invitation{
		InviterId:    creator.ID,
		InviteeId:    user1.ID,
		ResourceType: models.ResourceTypeChallenge,
		ResourceID:   created.ID,
		Status:       models.StatusPending,
	}
	config.DB.Create(&invitation)

	// 2. Accept invitation (this calls addUserToChallenge internally)
	err := services.AcceptInvitation(invitation.ID, user1.ID)
	assert.NoError(t, err)

	// 3. Verify user1 was added to challenge
	fetched, _ = services.GetChallengeByID(created.ID)
	assert.Len(t, fetched.Users, 2, "Challenge should have 2 users (creator + user1)")

	userIDs := make(map[uint]bool)
	for _, u := range fetched.Users {
		userIDs[u.ID] = true
	}
	assert.True(t, userIDs[creator.ID], "Creator should be in challenge")
	assert.True(t, userIDs[user1.ID], "User1 should be in challenge after accepting invitation")

	// 4. Create another invitation for user2 and accept it
	invitation2 := models.Invitation{
		InviterId:    creator.ID,
		InviteeId:    user2.ID,
		ResourceType: models.ResourceTypeChallenge,
		ResourceID:   created.ID,
		Status:       models.StatusPending,
	}
	config.DB.Create(&invitation2)

	err = services.AcceptInvitation(invitation2.ID, user2.ID)
	assert.NoError(t, err)

	// 5. Verify both users are in challenge
	fetched, _ = services.GetChallengeByID(created.ID)
	assert.Len(t, fetched.Users, 3, "Challenge should have 3 users")

	userIDs = make(map[uint]bool)
	for _, u := range fetched.Users {
		userIDs[u.ID] = true
	}
	assert.True(t, userIDs[creator.ID])
	assert.True(t, userIDs[user1.ID])
	assert.True(t, userIDs[user2.ID])

	// 6. Test error case: Try to accept invitation for non-existent challenge
	nonExistentInvitation := models.Invitation{
		InviterId:    creator.ID,
		InviteeId:    user1.ID,
		ResourceType: models.ResourceTypeChallenge,
		ResourceID:   99999, // Non-existent challenge
		Status:       models.StatusPending,
	}
	config.DB.Create(&nonExistentInvitation)

	err = services.AcceptInvitation(nonExistentInvitation.ID, user1.ID)
	assert.Error(t, err, "Should error when challenge doesn't exist")

	// 7. Test error case: Try to accept invitation for non-existent user
	nonExistentUserInvitation := models.Invitation{
		InviterId:    creator.ID,
		InviteeId:    99999, // Non-existent user
		ResourceType: models.ResourceTypeChallenge,
		ResourceID:   created.ID,
		Status:       models.StatusPending,
	}
	config.DB.Create(&nonExistentUserInvitation)

	err = services.AcceptInvitation(nonExistentUserInvitation.ID, 99999)
	assert.Error(t, err, "Should error when user doesn't exist")
}
