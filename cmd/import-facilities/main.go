package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"server/common/config"
	"server/common/models"
	"server/common/models/types"
)

// FacilityJSON matches the structure in facilities.json
type FacilityJSON struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	DetailedName string       `json:"detailedName"`
	Address      string       `json:"address"`
	Website      string       `json:"website"`
	Email        string       `json:"email"`
	FacilityType string       `json:"facilityType"`
	Indoor       bool         `json:"indoor"`
	Notes        string       `json:"notes"`
	Location     LocationJSON `json:"location"`
}

type LocationJSON struct {
	Address    string  `json:"address"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	PostalCode string  `json:"postal_code"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
}

func main() {
	log.Println("🏟️  Starting facility import...")

	config.LoadConfig()
	config.ConnectDatabase()

	if err := config.RunAtlasMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	path := "facilities.json"
	if p := os.Getenv("FACILITIES_JSON"); p != "" {
		path = p
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", path, err)
	}

	var facilities []FacilityJSON
	if err := json.Unmarshal(data, &facilities); err != nil {
		log.Fatalf("Failed to parse facilities JSON: %v", err)
	}

	log.Printf("📂 Loaded %d facilities from %s", len(facilities), path)

	batchSize := 500
	imported := 0
	skipped := 0

	for i := 0; i < len(facilities); i += batchSize {
		end := i + batchSize
		if end > len(facilities) {
			end = len(facilities)
		}
		batch := facilities[i:end]

		for _, f := range batch {
			address := f.Address
			if f.Location.Address != "" {
				address = f.Location.Address
			}
			facility := models.Facility{
				ExternalID:   f.ID,
				Name:         f.Name,
				DetailedName: f.DetailedName,
				Address:      address,
				Coordinates:  types.Point{Lat: f.Location.Latitude, Lon: f.Location.Longitude},
				PostalCode:   f.Location.PostalCode,
				City:         f.Location.City,
				Country:      f.Location.Country,
				FacilityType: f.FacilityType,
				Indoor:       f.Indoor,
			}

			if f.Website != "" {
				facility.Website = &f.Website
			}
			if f.Email != "" {
				facility.Email = &f.Email
			}
			if f.Notes != "" {
				facility.Notes = &f.Notes
			}

			result := config.DB.Create(&facility)
			if result.Error != nil {
				// Skip duplicates (unique constraint on external_id)
				if isDuplicateError(result.Error) {
					skipped++
					continue
				}
				log.Fatalf("Failed to create facility %s: %v", f.ID, result.Error)
			}
			imported++
		}

		log.Printf("  Progress: %d/%d (imported: %d, skipped: %d)", end, len(facilities), imported, skipped)
	}

	log.Printf("✅ Facility import complete: %d imported, %d skipped (duplicates)", imported, skipped)
}

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "unique constraint")
}
