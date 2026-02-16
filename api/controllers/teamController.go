package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"server/api/controllers/helpers"
	"server/common/appError"
	"server/common/dto"
	"server/common/middleware"
	"server/common/models"
	"server/common/services"
	"server/common/validator"
)

// --- GET ---
func GetTeam(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	teamModel, err := services.GetTeamByID(id, user.ID)
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
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	teamsModel, err := services.GetTeams(user.ID)
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
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	teamMembers, err := services.GetTeamsByUserId(id, user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Extract teams from team members
	response := make([]dto.TeamResponseDto, len(teamMembers))
	for i, tm := range teamMembers {
		response[i] = dto.ToTeamResponseDto(tm.Team)
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
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	teamMembers, err := services.GetTeamsByUserId(user.ID, user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Extract teams from team members
	response := make([]dto.TeamResponseDto, len(teamMembers))
	for i, tm := range teamMembers {
		response[i] = dto.ToTeamResponseDto(tm.Team)
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
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	req := dto.TeamCreateDto{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	if err := validator.V.Struct(req); err != nil {
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
	// Create team conversation with initial members
	memberIDs := make([]uint, len(createdModel.Users))
	for i, u := range createdModel.Users {
		memberIDs[i] = u.UserID
	}
	if err := services.SyncTeamConversationMembers(createdModel.ID, memberIDs); err != nil {
		// Log error but don't fail the request
		// Team conversation can be created later
		slog.Warn("Failed to create team conversation for team",
			slog.Int("team_id", int(createdModel.ID)),
			slog.Any("error", err),
		)

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createdDto)
		if err != nil {
			appError.HandleError(w, err)
			return
		}
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
	if err := validator.V.Struct(req); err != nil {
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
func SoftDeleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.SoftDeleteTeam(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Write 204
	w.WriteHeader(http.StatusNoContent)
}

func RemoveUserFromTeam(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)

	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	teamId, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	rmvUserId, err := helpers.GetParamIdDynamic(r, "rmvUserId")
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.RemoveUserFromTeam(*user, teamId, rmvUserId)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func LeaveTeam(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)

	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	teamId, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.LeaveTeam(*user, teamId)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
