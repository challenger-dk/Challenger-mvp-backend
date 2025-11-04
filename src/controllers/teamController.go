package controllers

import (
	"encoding/json"
	"net/http"
	"server/controllers/helpers"
	"server/dto"
	"server/services"
)

func GetTeam(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamId(w, r)
	if id == 0 {
		return
	}

	team, err := services.GetTeamByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(team)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetTeams(w http.ResponseWriter, r *http.Request) {
	users, err := services.GetTeams()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CreateTeam(w http.ResponseWriter, r *http.Request) {
	req := dto.TeamCreateDto{}

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created_team, err := services.CreateTeam(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(created_team)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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

	err = services.UpdateTeam(id, req)

	// Maybe this should be changed to something else
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Successfull update
	w.WriteHeader(http.StatusNoContent)
}

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
