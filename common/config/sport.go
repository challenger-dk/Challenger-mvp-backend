package config

import (
	"log/slog"
	"os"
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

	loadSportsCache()
	return nil
}

func loadSportsCache() {
	var sports []models.Sport

	err := DB.Model(&models.Sport{}).
		Select("name").
		Find(&sports).
		Error

	if err != nil {
		slog.Error("Failed to load sports cache", "error", err)
		os.Exit(1)
	}

	SportsCache = make(map[string]bool)
	for _, sport := range sports {
		SportsCache[sport.Name] = true
	}

	slog.Info("✅ Loaded sports into cache", "count", len(sports))
}
