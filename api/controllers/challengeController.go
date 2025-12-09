package controllers

import (
	"encoding/json"
	"net/http"
	"server/api/controllers/helpers"
	"server/api/middleware"
	"server/common/appError"
	"server/common/dto"
	"server/common/models"
	"server/common/services"
	"server/common/validator"
)

func GetChallenge(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	chalModel, err := services.GetChallengeByID(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	resp := dto.ToChallengeResponseDto(chalModel)

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		appError.HandleError(w, err)
	}
}

func GetChallenges(w http.ResponseWriter, r *http.Request) {
	challengesModel, err := services.GetChallenges()
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to response DTOs
	response := make([]dto.ChallengeResponseDto, len(challengesModel))
	for i, c := range challengesModel {
		response[i] = dto.ToChallengeResponseDto(c)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
	}
}

func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)

	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	req := dto.ChallengeCreateDto{}

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	challengeModel := dto.ChallengeCreateDtoToModel(req)
	// Set the creator ID to the authenticated user
	challengeModel.CreatorID = user.ID

	createdModel, err := services.CreateChallenge(challengeModel, req.Users)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	createdDto := dto.ToChallengeResponseDto(createdModel)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdDto)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func UpdateChallenge(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	req := dto.ChallengeCreateDto{}

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

	err = services.UpdateChallenge(id, dto.ChallengeCreateDtoToModel(req))

	// Maybe this should be changed to something else
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Successfull update
	w.WriteHeader(http.StatusNoContent)
}

func JoinChallenge(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	err = services.JoinChallenge(id, user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func LeaveChallenge(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	err = services.LeaveChallenge(id, user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteChallenge(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.DeleteChallenge(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
