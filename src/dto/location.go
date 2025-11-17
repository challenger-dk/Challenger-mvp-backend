package dto

import (
	"server/models"
	"server/models/types"
)

type LocationCreateDto struct {
	Address   string  `json:"address"   validate:"required"`
	Latitude  float64 `json:"latitude"  validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
}

type LocationResponseDto struct {
	ID        uint    `json:"id"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func LocationCreateDtoToModel(l LocationCreateDto) models.Location {
	return models.Location{
		Address: l.Address,
		Coordinates: types.Point{
			Lon: l.Longitude,
			Lat: l.Latitude,
		},
	}
}

func ToLocationResponseDto(l models.Location) LocationResponseDto {
	return LocationResponseDto{
		ID:        l.ID,
		Address:   l.Address,
		Latitude:  l.Coordinates.Lat,
		Longitude: l.Coordinates.Lon,
	}
}
