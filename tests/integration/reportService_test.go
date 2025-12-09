package integration

import (
	"server/common/config"
	"server/common/dto"
	"server/common/models"
	"server/common/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportService_CreateReport(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	// Setup: Create a reporter (User A) and a target (User B)
	reporter, err := services.CreateUser(models.User{
		Email:     "reporter@test.com",
		FirstName: "Reporter",
		LastName:  "Guy",
	}, "password")
	assert.NoError(t, err)

	targetUser, err := services.CreateUser(models.User{
		Email:     "target@test.com",
		FirstName: "Bad",
		LastName:  "Actor",
	}, "password")
	assert.NoError(t, err)

	t.Run("Report a User Success", func(t *testing.T) {
		req := dto.ReportCreateDto{
			TargetID:   targetUser.ID,
			TargetType: "USER",
			Reason:     "Harassment",
			Comment:    "He was mean in chat",
		}

		err := services.CreateReport(reporter.ID, req)
		assert.NoError(t, err)

		// Verify in DB
		var report models.Report
		err = config.DB.Where("reporter_id = ? AND target_id = ? AND target_type = ?", reporter.ID, targetUser.ID, "USER").First(&report).Error
		assert.NoError(t, err)
		assert.Equal(t, "Harassment", report.Reason)
		assert.Equal(t, "He was mean in chat", report.Comment)
		assert.Equal(t, "PENDING", report.Status)
	})

	t.Run("Report a Team Success", func(t *testing.T) {
		// Create a team to report
		team := models.Team{
			Name:      "Offensive Team Name",
			CreatorID: targetUser.ID,
		}
		config.DB.Create(&team)

		req := dto.ReportCreateDto{
			TargetID:   team.ID,
			TargetType: "TEAM",
			Reason:     "Inappropriate Name",
		}

		err := services.CreateReport(reporter.ID, req)
		assert.NoError(t, err)

		// Verify in DB
		var report models.Report
		err = config.DB.Where("target_id = ? AND target_type = ?", team.ID, "TEAM").First(&report).Error
		assert.NoError(t, err)
		assert.Equal(t, "Inappropriate Name", report.Reason)
	})

	t.Run("Report a Challenge Success", func(t *testing.T) {
		// Create a challenge
		challenge := models.Challenge{
			Name:      "Dangerous Event",
			CreatorID: targetUser.ID,
			Sport:     "Football",
		}
		config.DB.Create(&challenge)

		req := dto.ReportCreateDto{
			TargetID:   challenge.ID,
			TargetType: "CHALLENGE",
			Reason:     "Safety Concern",
		}

		err := services.CreateReport(reporter.ID, req)
		assert.NoError(t, err)

		// Verify
		var count int64
		config.DB.Model(&models.Report{}).Where("target_type = ?", "CHALLENGE").Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Report with Invalid Reporter ID (Fail)", func(t *testing.T) {
		// Attempt to create a report from a non-existent user ID
		// This should fail due to Foreign Key constraint on ReporterID
		req := dto.ReportCreateDto{
			TargetID:   targetUser.ID,
			TargetType: "USER",
			Reason:     "Spam",
		}

		nonExistentUserID := uint(99999)
		err := services.CreateReport(nonExistentUserID, req)

		assert.Error(t, err)
		// Verify it's a constraint violation
		assert.Contains(t, err.Error(), "violates foreign key constraint")
	})
}
