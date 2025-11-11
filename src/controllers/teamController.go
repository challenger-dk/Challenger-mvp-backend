package controllers

import (
	"encoding/json"
	"net/http"
	"server/controllers/helpers"
	"server/dto"
	"server/services"
)

// --- GET ---
func GetTeam(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamId(w, r)
	if id == 0 {
		return
	}

	teamModel, err := services.GetTeamByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := dto.ToTeamResponseDto(teamModel)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetTeams(w http.ResponseWriter, r *http.Request) {
	teamsModel, err := services.GetTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	response := make([]dto.TeamResponseDto, len(teamsModel))
	for i, t := range teamsModel {
		response[i] = dto.ToTeamResponseDto(t)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetTeamsByUserId(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamId(w, r)
	if id == 0 {
		return
	}

	teamsModel, err := services.GetTeamsByUserId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]dto.TeamResponseDto, len(teamsModel))
	for i, t := range teamsModel {
		response[i] = dto.ToTeamResponseDto(t)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// --- POST ---
func CreateTeam(w http.ResponseWriter, r *http.Request) {
	req := dto.TeamCreateDto{}

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert DTO -> model
	modelTeam := dto.TeamCreateDtoToModel(req)

	createdModel, err := services.CreateTeam(modelTeam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert model -> DTO for response
	createdDto := dto.ToTeamResponseDto(createdModel)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Depricated: Use Invitation system instead
/*
func AddUserToTeam(w http.ResponseWriter, r *http.Request) {
	teamId := helpers.GetParamId(w, r)
	if teamId == 0 {
		return
	}

	req := struct {
		UserId uint `json:"user_id"`
	}{}

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = services.AddUserToTeam(req.UserId, teamId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
*/

// --- PUT ---
func UpdateTeam(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamId(w, r)
	if id == 0 {
		return
	}

	req := dto.TeamCreateDto{}

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	modelTeam := dto.TeamCreateDtoToModel(req)

	err = services.UpdateTeam(id, modelTeam)

	// Maybe this should be changed to something else
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Successfull update
	w.WriteHeader(http.StatusNoContent)
}

// --- DELETE ---
func DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamId(w, r)
	if id == 0 {
		return
	}

	err := services.DeleteTeam(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
