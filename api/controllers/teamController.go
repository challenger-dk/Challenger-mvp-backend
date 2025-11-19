package controllers

import (
	"encoding/json"
	"net/http"
	"server/api/controllers/helpers"
	"server/api/middleware"
	"server/api/services"
	"server/common/appError"
	"server/common/dto"
	"server/common/models"
)

// --- GET ---
func GetTeam(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	teamModel, err := services.GetTeamByID(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	resp := dto.ToTeamResponseDto(teamModel)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetTeams(w http.ResponseWriter, r *http.Request) {
	teamsModel, err := services.GetTeams()
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to response DTOs
	response := make([]dto.TeamResponseDto, len(teamsModel))
	for i, t := range teamsModel {
		response[i] = dto.ToTeamResponseDto(t)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
	}
}

func GetTeamsByUserId(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	teamsModel, err := services.GetTeamsByUserId(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	response := make([]dto.TeamResponseDto, len(teamsModel))
	for i, t := range teamsModel {
		response[i] = dto.ToTeamResponseDto(t)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetCurrentUserTeams(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	teamsModel, err := services.GetTeamsByUserId(user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	response := make([]dto.TeamResponseDto, len(teamsModel))
	for i, t := range teamsModel {
		response[i] = dto.ToTeamResponseDto(t)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

// --- POST ---
func CreateTeam(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	req := dto.TeamCreateDto{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	if err := validate.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	modelTeam := dto.TeamCreateDtoToModel(req)

	// Set the creator ID to the authenticated user
	modelTeam.CreatorID = user.ID

	createdModel, err := services.CreateTeam(modelTeam)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	createdDto := dto.ToTeamResponseDto(createdModel)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdDto)
	if err != nil {
		appError.HandleError(w, err)
		return
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
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	req := dto.TeamUpdateDto{}

	// Decode request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	modelTeam := dto.TeamUpdateDtoToModel(req)

	err = services.UpdateTeam(id, modelTeam)

	// Maybe this should be changed to something else
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Successfull update
	w.WriteHeader(http.StatusNoContent)
}

// --- DELETE ---
func DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.DeleteTeam(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}
