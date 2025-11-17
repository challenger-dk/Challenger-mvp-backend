package controllers

import (
	"encoding/json"
	"net/http"
	"server/appError"
	"server/dto"
	"server/services"
)

func GetSports(w http.ResponseWriter, r *http.Request) {
	sports, err := services.GetAllSports()
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to response DTOs
	response := make([]dto.SportResponseDto, len(sports))
	for i, sport := range sports {
		response[i] = dto.ToSportResponseDto(sport)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}
