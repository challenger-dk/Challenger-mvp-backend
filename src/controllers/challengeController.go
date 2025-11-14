package controllers

import (
	"encoding/json"
	"net/http"
	"server/controllers/helpers"
	"server/dto"
	"server/services"
)

func GetChallenge(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamId(w, r)
	if id == 0 {
		return
	}

	chalModel, err := services.GetChallengeByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := dto.ToChallengeResponseDto(chalModel)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetChallenges(w http.ResponseWriter, r *http.Request) {
	challengesModel, err := services.GetChallenges()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	response := make([]dto.ChallengeResponseDto, len(challengesModel))
	for i, c := range challengesModel {
		response[i] = dto.ToChallengeResponseDto(c)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CreateChallenge(w http.ResponseWriter, r *http.Request) {
	req := dto.ChallengeCreateDto{}

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdModel, err := services.CreateChallenge(dto.ChallengeCreateDtoToModel(req))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	createdDto := dto.ToChallengeResponseDto(createdModel)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateChallenge(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamId(w, r)
	if id == 0 {
		return
	}

	req := dto.ChallengeCreateDto{}

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = services.UpdateChallenge(id, dto.ChallengeCreateDtoToModel(req))

	// Maybe this should be changed to something else
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Successfull update
	w.WriteHeader(http.StatusNoContent)
}

func DeleteChallenge(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamId(w, r)
	if id == 0 {
		return
	}

	err := services.DeleteChallenge(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
