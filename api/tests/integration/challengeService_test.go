package integration

import (
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
