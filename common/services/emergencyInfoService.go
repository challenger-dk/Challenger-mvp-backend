package services

import (
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

func CreateEmergencyContact(user models.User, emergencyInfo models.EmergencyInfo) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Check if user exists
		var existingUser models.User
		if err := tx.First(&existingUser, user.ID).Error; err != nil {
			return err
		}

		// Set the UserID for the emergency contact
		emergencyInfoModel := emergencyInfo
		emergencyInfoModel.UserID = user.ID

		// Create the emergency contact
		if err := tx.Create(&emergencyInfoModel).Error; err != nil {
			return err
		}

		return nil
	})
}

func DeleteEmergencyContact(user models.User, emergencyInfoID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Check if contact exists and belongs to user
		var emergencyInfo models.EmergencyInfo
		if err := tx.First(&emergencyInfo, emergencyInfoID).Error; err != nil {
			return err
		}

		// Verify ownership
		if emergencyInfo.UserID != user.ID {
			return gorm.ErrRecordNotFound
		}

		// Delete emergency contact
		if err := tx.Delete(&emergencyInfo).Error; err != nil {
			return err
		}

		return nil
	})
}

func UpdateEmergencyContact(user models.User, emergencyInfo models.EmergencyInfo, emergencyInfoID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Check if contact exists and belongs to user
		var existingEmergencyInfo models.EmergencyInfo
		if err := tx.First(&existingEmergencyInfo, emergencyInfoID).Error; err != nil {
			return err
		}

		// Verify ownership
		if existingEmergencyInfo.UserID != user.ID {
			return gorm.ErrRecordNotFound
		}

		// Update fields
		existingEmergencyInfo.Name = emergencyInfo.Name
		existingEmergencyInfo.PhoneNumber = emergencyInfo.PhoneNumber
		existingEmergencyInfo.Relationship = emergencyInfo.Relationship

		// Save changes
		if err := tx.Save(&existingEmergencyInfo).Error; err != nil {
			return err
		}

		return nil
	})
}
