package integration

import (
	"server/common/config"
	"server/common/models"
	"server/common/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSportService_GetAllSports(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Sports are already seeded in setupTest via config.SeedSports()
	// So we can test retrieving them

	// 1. Get all sports
	sports, err := services.GetAllSports()
	assert.NoError(t, err)
	assert.NotNil(t, sports)

	// 2. Verify we got the expected number of sports
	// GetAllowedSports() returns 22 sports
	expectedSports := models.GetAllowedSports()
	assert.GreaterOrEqual(t, len(sports), len(expectedSports), "Should have at least the seeded sports")

	// 3. Verify sports have valid structure
	for _, sport := range sports {
		assert.NotZero(t, sport.ID, "Sport should have an ID")
		assert.NotEmpty(t, sport.Name, "Sport should have a name")
		assert.False(t, sport.CreatedAt.IsZero(), "Sport should have CreatedAt timestamp")
	}

	// 4. Verify all expected sports are present
	sportNames := make(map[string]bool)
	for _, sport := range sports {
		sportNames[sport.Name] = true
	}

	for _, expectedSportName := range expectedSports {
		assert.True(t, sportNames[expectedSportName], "Expected sport '%s' should be present", expectedSportName)
	}

	// 5. Verify no duplicate sports (unique constraint)
	sportIDMap := make(map[uint]bool)
	sportNameMap := make(map[string]bool)
	for _, sport := range sports {
		assert.False(t, sportIDMap[sport.ID], "Sport IDs should be unique, found duplicate ID: %d", sport.ID)
		assert.False(t, sportNameMap[sport.Name], "Sport names should be unique, found duplicate name: %s", sport.Name)
		sportIDMap[sport.ID] = true
		sportNameMap[sport.Name] = true
	}
}

func TestSportService_GetAllSports_EmptyDatabase(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Clear all sports from database
	config.DB.Exec("TRUNCATE TABLE sports RESTART IDENTITY CASCADE")

	// Get all sports from empty database
	sports, err := services.GetAllSports()
	assert.NoError(t, err)
	assert.NotNil(t, sports)
	assert.Empty(t, sports, "Should return empty slice when no sports exist")
}

func TestSportService_GetAllSports_AfterAddingNewSport(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Get initial count
	initialSports, err := services.GetAllSports()
	assert.NoError(t, err)
	initialCount := len(initialSports)

	// Add a new sport directly to database (simulating a new sport being added)
	newSport := models.Sport{Name: "TestSport"}
	err = config.DB.Create(&newSport).Error
	assert.NoError(t, err)

	// Get all sports again
	allSports, err := services.GetAllSports()
	assert.NoError(t, err)
	assert.Equal(t, initialCount+1, len(allSports), "Should have one more sport after adding")

	// Verify the new sport is in the results
	found := false
	for _, sport := range allSports {
		if sport.Name == "TestSport" {
			found = true
			assert.NotZero(t, sport.ID)
			break
		}
	}
	assert.True(t, found, "New sport should be in the results")
}
