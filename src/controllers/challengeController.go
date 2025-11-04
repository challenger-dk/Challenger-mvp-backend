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

	chal, err := services.GetChallengeByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(chal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetChallenges(w http.ResponseWriter, r *http.Request) {
	challenges, err := services.GetChallenges()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(challenges)
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

	created_chal, err := services.CreateChallenge(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(created_chal)
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

	err = services.UpdateChallenge(id, req)

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
