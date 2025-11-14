package controllers

import (
	"encoding/json"
	"net/http"
	"server/dto"
	"server/services"
)

func GetSports(w http.ResponseWriter, r *http.Request) {
	sports, err := services.GetAllSports()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	response := make([]dto.SportResponseDto, len(sports))
	for i, sport := range sports {
		response[i] = dto.ToSportResponseDto(sport)
	}

	json.NewEncoder(w).Encode(response)
}
