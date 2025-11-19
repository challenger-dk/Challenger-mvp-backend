package services

import (
	"server/common/models"

	"gorm.io/gorm"
)

// FindOrCreateLocation finds an existing location by coordinates or creates a new one.
// It is designed to be run inside a transaction.
func FindOrCreateLocation(tx *gorm.DB, locationModel models.Location) (*models.Location, error) {
	var location models.Location

	err := tx.Where(models.Location{Coordinates: locationModel.Coordinates}).
		FirstOrCreate(&location, locationModel).
		Error

	if err != nil {
		return nil, err
	}

	return &location, nil
}
