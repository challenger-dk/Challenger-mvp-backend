package controllers

import (
	"encoding/json"
	"net/http"
	"server/common/appError"
	"server/common/dto"
	"server/common/services"
)

func GetSports(w http.ResponseWriter, r *http.Request) {
	sports, err := services.GetAllSports()
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to response DTOs
	response := make([]dto.SportDto, len(sports))
	for i, sport := range sports {
		response[i] = dto.ToSportDto(sport)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}
