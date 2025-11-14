package config // or database

import (
	"log"
	"server/models"
)

// SportsCache holds the allowed sports for quick validation
var SportsCache = make(map[string]bool)

// Call this ONCE from main.go
func LoadSportsCache() {
	var sports []models.Sport

	// Use the DB connection we already have
	// Assumes your sport model is models.Sport and table is "sports"
	if err := DB.Model(&models.Sport{}).Select("name").Find(&sports).Error; err != nil {
		log.Fatalf("Failed to load sports cache: %v", err)
	}

	// Clear the cache and reload
	SportsCache = make(map[string]bool)
	for _, sport := range sports {
		SportsCache[sport.Name] = true
	}

	log.Printf("✅ Loaded %d sports into cache", len(sports))
}

// Should be called before LoadSportsCache
func SeedSports() error {
	allowedSports := models.GetAllowedSports()

	for _, sportName := range allowedSports {
		var sport models.Sport
		if err := DB.Where("name = ?", sportName).FirstOrCreate(&sport, models.Sport{Name: sportName}).Error; err != nil {
			return err
		}
	}

	return nil
}
