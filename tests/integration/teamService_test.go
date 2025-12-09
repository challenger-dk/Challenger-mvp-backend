package integration

import (
	"server/common/appError"
	"server/common/models"
	"server/common/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamService_CRUD(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	creator, _ := services.CreateUser(models.User{Email: "creator@team.com", FirstName: "C", LastName: "C"}, "pw")

	// 1. Create Team
	teamModel := models.Team{
		Name:      "Test Team",
		CreatorID: creator.ID,
		Location: &models.Location{
			Address:     "Test St",
			City:        "Test City",
			Country:     "DK",
			PostalCode:  "1000",
			Coordinates: models.Point{Lat: 55.0, Lon: 12.0},
		},
	}

	createdTeam, err := services.CreateTeam(teamModel)
	assert.NoError(t, err)
	assert.NotZero(t, createdTeam.ID)
	assert.Equal(t, creator.ID, createdTeam.CreatorID)
	assert.Len(t, createdTeam.Users, 1) // Creator auto-added

	// 2. Get Team By ID
	fetched, err := services.GetTeamByID(createdTeam.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Test Team", fetched.Name)

	// 3. Update Team
	err = services.UpdateTeam(createdTeam.ID, models.Team{Name: "Updated Team"})
	assert.NoError(t, err)

	updated, _ := services.GetTeamByID(createdTeam.ID)
	assert.Equal(t, "Updated Team", updated.Name)

	// 4. Delete Team
	err = services.DeleteTeam(createdTeam.ID)
	assert.NoError(t, err)

	_, err = services.GetTeamByID(createdTeam.ID)
	assert.Error(t, err)
}

func TestTeamService_List(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	u, _ := services.CreateUser(models.User{Email: "u@team.com", FirstName: "U", LastName: "U"}, "pw")
	services.CreateTeam(models.Team{Name: "T1", CreatorID: u.ID})
	services.CreateTeam(models.Team{Name: "T2", CreatorID: u.ID})

	// Get All
	teams, err := services.GetTeams()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(teams), 2)

	// Get By User
	userTeams, err := services.GetTeamsByUserId(u.ID)
	assert.NoError(t, err)
	assert.Len(t, userTeams, 2)
}

func TestTeamService_Membership(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	creator, _ := services.CreateUser(models.User{Email: "c@t.com", FirstName: "C", LastName: "C"}, "pw")
	member, _ := services.CreateUser(models.User{Email: "m@t.com", FirstName: "M", LastName: "M"}, "pw")
	outsider, _ := services.CreateUser(models.User{Email: "o@t.com", FirstName: "O", LastName: "O"}, "pw")

	team, _ := services.CreateTeam(models.Team{Name: "T", CreatorID: creator.ID})

	// Manually add member for test setup
	invitation := &models.Invitation{
		InviterId:    creator.ID,
		InviteeId:    member.ID,
		ResourceType: models.ResourceTypeTeam,
		ResourceID:   team.ID,
	}
	services.SendInvitation(invitation)
	// Re-fetch to get ID
	invites, _ := services.GetInvitationsByUserId(member.ID)
	services.AcceptInvitation(invites[0].ID, member.ID)

	// 1. Remove User (Unauthorized - Outsider tries to remove Member)
	err := services.RemoveUserFromTeam(*outsider, team.ID, member.ID)
	assert.ErrorIs(t, err, appError.ErrUnauthorized)

	// 2. Remove User (Authorized - Creator removes Member)
	err = services.RemoveUserFromTeam(*creator, team.ID, member.ID)
	assert.NoError(t, err)

	// Verify removal
	tAfterRemove, _ := services.GetTeamByID(team.ID)
	assert.Len(t, tAfterRemove.Users, 1) // Only creator left

	// 3. Leave Team
	err = services.LeaveTeam(*creator, team.ID)
	assert.NoError(t, err)

	tAfterLeave, _ := services.GetTeamByID(team.ID)
	assert.Len(t, tAfterLeave.Users, 0)
}
