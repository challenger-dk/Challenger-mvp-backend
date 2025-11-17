package dto

import (
	"server/models"
)

type SportDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func ToSportDto(sport models.Sport) SportDto {
	return SportDto{
		ID:   sport.ID,
		Name: sport.Name,
	}
}
