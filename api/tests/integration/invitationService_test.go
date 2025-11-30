package integration

import (
	"server/api/services"
	"server/common/appError"
	"server/common/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvitationService_TeamFlow(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	inviter, _ := services.CreateUser("inv@t.com", "pw", "I", "I", nil)
	invitee, _ := services.CreateUser("rec@t.com", "pw", "R", "R", nil)
	team, _ := services.CreateTeam(models.Team{Name: "T", CreatorID: inviter.ID})

	// 1. Send Invite
	inv := &models.Invitation{
		InviterId:    inviter.ID,
		InviteeId:    invitee.ID,
		ResourceType: models.ResourceTypeTeam,
		ResourceID:   team.ID,
	}
	err := services.SendInvitation(inv)
	assert.NoError(t, err)

	// 2. Send Duplicate (Should return Pending Error)
	err = services.SendInvitation(inv)
	assert.ErrorIs(t, err, appError.ErrInvitationPending)

	// 3. Accept Invite
	// Retrieve created ID
	invites, _ := services.GetInvitationsByUserId(invitee.ID)
	assert.NotEmpty(t, invites)

	err = services.AcceptInvitation(invites[0].ID, invitee.ID)
	assert.NoError(t, err)

	// 4. Verify Accepted Status
	// Note: AcceptInvitation updates the row.
	// But `SendInvitation` logic says if StatusAccepted -> Return ErrInvitationAccepted
	err = services.SendInvitation(inv)
	assert.ErrorIs(t, err, appError.ErrInvitationAccepted)

	// 5. Verify Team Member
	updatedTeam, _ := services.GetTeamByID(team.ID)
	assert.Len(t, updatedTeam.Users, 2)
}

func TestInvitationService_FriendRequest(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	u1, _ := services.CreateUser("u1@f.com", "pw", "1", "1", nil)
	u2, _ := services.CreateUser("u2@f.com", "pw", "2", "2", nil)

	// 1. Send Friend Req
	inv := &models.Invitation{
		InviterId:    u1.ID,
		InviteeId:    u2.ID,
		ResourceType: models.ResourceTypeFriend,
		ResourceID:   0, // Not used
	}
	err := services.SendInvitation(inv)
	assert.NoError(t, err)

	// 2. Get Invites
	list, err := services.GetInvitationsByUserId(u2.ID)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// 3. Accept
	err = services.AcceptInvitation(list[0].ID, u2.ID)
	assert.NoError(t, err)
}

func TestInvitationService_Errors(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	u1, _ := services.CreateUser("e1@e.com", "pw", "1", "1", nil)
	u2, _ := services.CreateUser("e2@e.com", "pw", "2", "2", nil)

	// 1. Invite Self
	inv := &models.Invitation{InviterId: u1.ID, InviteeId: u1.ID, ResourceType: models.ResourceTypeFriend}
	err := services.SendInvitation(inv)
	assert.ErrorIs(t, err, appError.ErrInviteSameUser)

	// 2. Accept Unauthorized
	inv2 := &models.Invitation{InviterId: u1.ID, InviteeId: u2.ID, ResourceType: models.ResourceTypeFriend}
	services.SendInvitation(inv2)
	list, _ := services.GetInvitationsByUserId(u2.ID)

	// u1 tries to accept u2's invite
	err = services.AcceptInvitation(list[0].ID, u1.ID)
	assert.ErrorIs(t, err, appError.ErrUnauthorized)
}

func TestInvitationService_DeclineResend(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	u1, _ := services.CreateUser("d1@d.com", "pw", "1", "1", nil)
	u2, _ := services.CreateUser("d2@d.com", "pw", "2", "2", nil)

	inv := &models.Invitation{InviterId: u1.ID, InviteeId: u2.ID, ResourceType: models.ResourceTypeFriend}
	services.SendInvitation(inv)
	list, _ := services.GetInvitationsByUserId(u2.ID)

	// 1. Decline
	err := services.DeclineInvitation(list[0].ID, u2.ID)
	assert.NoError(t, err)

	// 2. Resend (Should succeed and reset status)
	err = services.SendInvitation(inv)
	assert.NoError(t, err)

	// Check status is Pending again
	list2, _ := services.GetInvitationsByUserId(u2.ID)
	assert.Equal(t, models.StatusPending, list2[0].Status)
}
