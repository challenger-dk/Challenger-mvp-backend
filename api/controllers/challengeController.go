package controllers

import (
	"encoding/json"
	"net/http"
	"server/api/controllers/helpers"
	"server/api/services"
	"server/common/appError"
	"server/common/dto"
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
	req := dto.ChallengeCreateDto{}

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	createdModel, err := services.CreateChallenge(dto.ChallengeCreateDtoToModel(req), req.Users)
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
	if err := validate.Struct(req); err != nil {
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
}
