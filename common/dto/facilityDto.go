package dto

import "server/common/models"

// FacilityResponseDto represents a facility in API responses.
// Matches the frontend facilities.json structure for compatibility.
type FacilityResponseDto struct {
	ID           uint    `json:"id"`
	ExternalID   string  `json:"external_id"`
	Name         string  `json:"name"`
	DetailedName string  `json:"detailed_name"`
	Address      string  `json:"address"`
	Website      *string `json:"website,omitempty"`
	Email        *string `json:"email,omitempty"`
	FacilityType string  `json:"facility_type"`
	Indoor       bool    `json:"indoor"`
	Notes        *string `json:"notes,omitempty"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	PostalCode   string  `json:"postal_code"`
	City         string  `json:"city"`
	Country      string  `json:"country"`
}

// ToFacilityResponseDto converts a Facility model to the response DTO.
func ToFacilityResponseDto(f models.Facility) FacilityResponseDto {
	return FacilityResponseDto{
		ID:           f.ID,
		ExternalID:   f.ExternalID,
		Name:         f.Name,
		DetailedName: f.DetailedName,
		Address:      f.Address,
		Website:      f.Website,
		Email:        f.Email,
		FacilityType: f.FacilityType,
		Indoor:       f.Indoor,
		Notes:        f.Notes,
		Latitude:     f.Coordinates.Lat,
		Longitude:    f.Coordinates.Lon,
		PostalCode:   f.PostalCode,
		City:         f.City,
		Country:      f.Country,
	}
}
