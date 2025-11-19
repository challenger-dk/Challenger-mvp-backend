package dto

import (
	"server/common/models"
	"server/common/models/types"
)

type LocationCreateDto struct {
	Address    string  `json:"address"   validate:"required"`
	Latitude   float64 `json:"latitude"  validate:"required,min=-90,max=90"`
	Longitude  float64 `json:"longitude" validate:"required,min=-180,max=180"`
	PostalCode string  `json:"postal_code" validate:"required"`
	City       string  `json:"city" validate:"required"`
	Country    string  `json:"country" validate:"required"`
}

type LocationResponseDto struct {
	ID         uint    `json:"id"`
	Address    string  `json:"address"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	PostalCode string  `json:"postal_code"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
}

func LocationCreateDtoToModel(l LocationCreateDto) models.Location {
	return models.Location{
		Address: l.Address,
		Coordinates: types.Point{
			Lon: l.Longitude,
			Lat: l.Latitude,
		},
		PostalCode: l.PostalCode,
		City:       l.City,
		Country:    l.Country,
	}
}

func ToLocationResponseDto(l models.Location) LocationResponseDto {
	return LocationResponseDto{
		ID:         l.ID,
		Address:    l.Address,
		Latitude:   l.Coordinates.Lat,
		Longitude:  l.Coordinates.Lon,
		PostalCode: l.PostalCode,
		City:       l.City,
		Country:    l.Country,
	}
}
