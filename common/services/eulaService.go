package services

import (
	"errors"

	"server/common/appError"
	"server/common/config"
	"server/common/models"

	"gorm.io/gorm"
)

// GetActiveEula returns the active EULA for the specified locale
func GetActiveEula(locale string) (*models.EulaVersion, error) {
	if locale == "" {
		locale = "da-DK" // Default locale
	}

	var eulaVersion models.EulaVersion
	err := config.DB.Where("locale = ? AND is_active = ?", locale, true).
		First(&eulaVersion).Error

	if err != nil {
		return nil, err
	}

	return &eulaVersion, nil
}

// GetUserEulaAcceptance checks if the user has accepted a specific EULA version
// Returns the acceptance record if found, or nil if not found (with no error)
func GetUserEulaAcceptance(userID uint, eulaVersionID uint) (*models.EulaAcceptance, error) {
	var acceptance models.EulaAcceptance
	err := config.DB.Where("user_id = ? AND eula_version_id = ?", userID, eulaVersionID).
		First(&acceptance).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error
		}
		return nil, err
	}

	return &acceptance, nil
}

// AcceptEula records the user's acceptance of a EULA version
func AcceptEula(userID uint, eulaVersionID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// Verify the EULA version exists and is active
		var eulaVersion models.EulaVersion
		err := tx.First(&eulaVersion, eulaVersionID).Error
		if err != nil {
			return err
		}

		// Verify it's the active version
		if !eulaVersion.IsActive {
			return appError.ErrEulaNotActive
		}

		// Check if user has already accepted this version (idempotent)
		var existing models.EulaAcceptance
		err = tx.Where("user_id = ? AND eula_version_id = ?", userID, eulaVersionID).
			First(&existing).Error

		if err == nil {
			// Already accepted - idempotent, just return success
			return nil
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// Actual error
			return err
		}

		// Create new acceptance record
		acceptance := models.EulaAcceptance{
			UserID:        userID,
			EulaVersionID: eulaVersionID,
		}

		if err := tx.Create(&acceptance).Error; err != nil {
			return err
		}

		return nil
	})
}

// HasUserAcceptedActiveEula checks if user has accepted the active EULA for their locale
// This is used by middleware for quick checks
func HasUserAcceptedActiveEula(userID uint, locale string) (bool, error) {
	activeEula, err := GetActiveEula(locale)
	if err != nil {
		return false, err
	}

	acceptance, err := GetUserEulaAcceptance(userID, activeEula.ID)
	if err != nil {
		return false, err
	}

	return acceptance != nil, nil
}
