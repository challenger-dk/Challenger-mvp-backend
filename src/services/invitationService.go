package services

import (
	"server/config"
	"server/models"
)

// --- GET ---
func GetInvitationsByUserId(id uint) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := config.DB.Preload("Team").Preload("FromUser").Where("to_user_id = ?", id).Find(&invitations).Error
	if err != nil {
		return nil, err
	}
	return invitations, nil
}
