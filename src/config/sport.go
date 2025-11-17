package config // or database

import (
	"log"
	"server/models"
)

// SportsCache holds the allowed sports for quick validation
var SportsCache = make(map[string]bool)

// Should be called before LoadSportsCache
func SeedSports() error {
	allowedSports := models.GetAllowedSports()

	for _, sportName := range allowedSports {
		var sport models.Sport

		err := DB.Where("name = ?", sportName).
			FirstOrCreate(&sport, models.Sport{Name: sportName}).
			Error

		if err != nil {
			return err
		}
	}

	loadSportsCache()

	return nil
}

// LoadSportsCache loads the allowed sports from the database into the SportsCache map
// Ensures the server autoamatically has the latest allowed sports in the database
func loadSportsCache() {
	var sports []models.Sport

	// Use the DB connection we already have
	// Assumes your sport model is models.Sport and table is "sports"
	err := DB.Model(&models.Sport{}).
		Select("name").
		Find(&sports).
		Error

	if err != nil {
		log.Fatalf("Failed to load sports cache: %v", err)
	}

	// Clear the cache and reload
	SportsCache = make(map[string]bool)
	for _, sport := range sports {
		SportsCache[sport.Name] = true
	}

	log.Printf("✅ Loaded %d sports into cache", len(sports))
}
