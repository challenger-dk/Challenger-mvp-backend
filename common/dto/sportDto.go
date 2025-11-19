package dto

import (
	"server/common/models"
)

type SportDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type SportResponseDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToSportDto(sport models.Sport) SportDto {
	return SportDto{
		ID:   sport.ID,
		Name: sport.Name,
	}
}

func ToSportResponseDto(sport models.Sport) SportResponseDto {
	return SportResponseDto{
		ID:   sport.ID,
		Name: sport.Name,
	}
}
