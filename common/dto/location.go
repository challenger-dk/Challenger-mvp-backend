package dto

import (
	"server/common/models"
	"server/common/models/types"
)

type LocationCreateDto struct {
	Type       string  `json:"type" validate:"sanitize,oneof=address facility"`
	Address    string  `json:"address"   validate:"sanitize,required"`
	Latitude   float64 `json:"latitude"  validate:"required,min=-90,max=90"`
	Longitude  float64 `json:"longitude" validate:"required,min=-180,max=180"`
	PostalCode string  `json:"postal_code" validate:"sanitize,required"`
	City       string  `json:"city" validate:"sanitize,required"`
	Country    string  `json:"country" validate:"sanitize,required"`
	// Facility-specific fields (optional)
	FacilityID   *string `json:"facility_id" validate:"omitempty,sanitize"`
	FacilityName *string `json:"facility_name" validate:"omitempty,sanitize"`
	DetailedName *string `json:"detailed_name" validate:"omitempty,sanitize"`
	Email        *string `json:"email" validate:"omitempty,sanitize"`
	Website      *string `json:"website" validate:"omitempty,sanitize"`
	FacilityType *string `json:"facility_type" validate:"omitempty,sanitize"`
	Indoor       *bool   `json:"indoor"`
	Notes        *string `json:"notes" validate:"omitempty,sanitize"`
}

type LocationResponseDto struct {
	ID         uint    `json:"id"`
	Type       string  `json:"type"`
	Address    string  `json:"address"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	PostalCode string  `json:"postal_code"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	// Facility-specific fields (optional)
	FacilityID   *string `json:"facility_id,omitempty"`
	FacilityName *string `json:"facility_name,omitempty"`
	DetailedName *string `json:"detailed_name,omitempty"`
	Email        *string `json:"email,omitempty"`
	Website      *string `json:"website,omitempty"`
	FacilityType *string `json:"facility_type,omitempty"`
	Indoor       *bool   `json:"indoor,omitempty"`
	Notes        *string `json:"notes,omitempty"`
}

func LocationCreateDtoToModel(l LocationCreateDto) models.Location {
	locationType := models.LocationTypeAddress
	if l.Type == "facility" {
		locationType = models.LocationTypeFacility
	}

	return models.Location{
		Type:    locationType,
		Address: l.Address,
		Coordinates: types.Point{
			Lon: l.Longitude,
			Lat: l.Latitude,
		},
		PostalCode:   l.PostalCode,
		City:         l.City,
		Country:      l.Country,
		FacilityID:   l.FacilityID,
		FacilityName: l.FacilityName,
		DetailedName: l.DetailedName,
		Email:        l.Email,
		Website:      l.Website,
		FacilityType: l.FacilityType,
		Indoor:       l.Indoor,
		Notes:        l.Notes,
	}
}

func ToLocationResponseDto(l models.Location) LocationResponseDto {
	return LocationResponseDto{
		ID:           l.ID,
		Type:         string(l.Type),
		Address:      l.Address,
		Latitude:     l.Coordinates.Lat,
		Longitude:    l.Coordinates.Lon,
		PostalCode:   l.PostalCode,
		City:         l.City,
		Country:      l.Country,
		FacilityID:   l.FacilityID,
		FacilityName: l.FacilityName,
		DetailedName: l.DetailedName,
		Email:        l.Email,
		Website:      l.Website,
		FacilityType: l.FacilityType,
		Indoor:       l.Indoor,
		Notes:        l.Notes,
	}
}
