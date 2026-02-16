package services

import (
	"server/common/config"
	"server/common/models"
)

// ListFacilities returns facilities with optional search and filters.
// Supports: search_query, city, limit, offset, min_lat, max_lat, min_lon, max_lon.
func ListFacilities(searchQuery string, city string, limit int, offset int, minLat, maxLat, minLon, maxLon float64, hasMinLat, hasMaxLat, hasMinLon, hasMaxLon bool) ([]models.Facility, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 2000 {
		limit = 2000
	}

	query := config.DB.Model(&models.Facility{})

	if searchQuery != "" {
		pattern := "%" + searchQuery + "%"
		query = query.Where(
			"name ILIKE ? OR detailed_name ILIKE ? OR city ILIKE ? OR address ILIKE ?",
			pattern, pattern, pattern, pattern,
		)
	}
	if city != "" {
		query = query.Where("city ILIKE ?", city)
	}

	// Geographic bounds (map viewport) - coordinates stored as PostGIS geography(Point,4326)
	// Use ST_MakeEnvelope(minLon, minLat, maxLon, maxLat, 4326) and && operator for bounding box
	if hasMinLat && hasMaxLat && hasMinLon && hasMaxLon {
		query = query.Where(
			"facilities.coordinates && ST_MakeEnvelope(?, ?, ?, ?, 4326)::geography",
			minLon, minLat, maxLon, maxLat,
		)
	}

	var facilities []models.Facility
	err := query.Limit(limit).Offset(offset).Find(&facilities).Error
	return facilities, err
}

// GetFacilityByID returns a facility by its database ID.
func GetFacilityByID(id uint) (models.Facility, error) {
	var facility models.Facility
	err := config.DB.First(&facility, id).Error
	return facility, err
}

// GetFacilityByExternalID returns a facility by its external ID (e.g. "facility-aabenraa-1").
func GetFacilityByExternalID(externalID string) (models.Facility, error) {
	var facility models.Facility
	err := config.DB.Where("external_id = ?", externalID).First(&facility).Error
	return facility, err
}
