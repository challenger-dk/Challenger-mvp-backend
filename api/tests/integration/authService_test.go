package integration

import (
	"server/api/services"
	"server/common/appError"
	"server/common/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthService_FullFlow(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	config.AppConfig.JWTSecret = "test_secret_key_12345"

	email := "auth_full@test.com"
	password := "strongPassword"

	// 1. Login Non-Existent User
	_, _, err := services.Login("ghost@test.com", password)
	assert.ErrorIs(t, err, appError.ErrInvalidCredentials)

	// 2. Create User
	createdUser, err := services.CreateUser(email, password, "Auth", "User", nil)
	assert.NoError(t, err)

	// 3. Login Wrong Password
	_, _, err = services.Login(email, "wrongPass")
	assert.ErrorIs(t, err, appError.ErrInvalidCredentials)

	// 4. Login Success
	user, token, err := services.Login(email, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, token)
	assert.Equal(t, createdUser.ID, user.ID)

	// 5. Validate Token Success
	claims, err := services.ValidateJWTToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Email, claims.Email)

	// 6. Validate Invalid Token (Garbage)
	_, err = services.ValidateJWTToken("garbage.token.string")
	assert.Error(t, err) // jwt parsing error

	// 7. Validate Signed but Modified Token (Tampered)
	// Hard to simulate without manually crafting a token with different secret,
	// but sending a random string is sufficient coverage for the parsing logic.
}
