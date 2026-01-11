package config

import (
	"log/slog"
	"server/common/models"
)

var SportsCache = make(map[string]bool)

func SeedSports() error {
	allowedSports := models.GetAllowedSports()

	for _, sportName := range allowedSports {
		var sport models.Sport
		// FirstOrCreate to avoid duplicates
		err := DB.Where("name = ?", sportName).
			FirstOrCreate(&sport, models.Sport{Name: sportName}).
			Error

		if err != nil {
			return err
		}
	}

	if err := LoadSportsCache(); err != nil {
		return err
	}
	return nil
}

// LoadSportsCache loads sports from the database into the cache
// Falls back to hardcoded list if database table doesn't exist yet
func LoadSportsCache() error {
	var sports []models.Sport

	err := DB.Model(&models.Sport{}).
		Select("name").
		Find(&sports).
		Error

	if err != nil {
		// If the table doesn't exist yet (migrations haven't run), use hardcoded list
		slog.Warn("Sports table not found, using hardcoded sports list", "error", err)
		allowedSports := models.GetAllowedSports()

		SportsCache = make(map[string]bool)
		for _, sportName := range allowedSports {
			SportsCache[sportName] = true
		}

		slog.Info("✅ Loaded sports into cache from hardcoded list", "count", len(allowedSports))
		return nil
	}

	SportsCache = make(map[string]bool)
	for _, sport := range sports {
		SportsCache[sport.Name] = true
	}

	slog.Info("✅ Loaded sports into cache from database", "count", len(sports))
	return nil
}
