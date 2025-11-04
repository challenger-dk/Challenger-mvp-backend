package controllers

import (
	"encoding/json"
	"net/http"
	"server/dto"
	"server/services"
	"strconv"
)

func GetTeam(w http.ResponseWriter, r *http.Request) {

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
	req := dto.RequestTeam{}

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created_user, err := services.CreateTeam(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(created_user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UpdateTeam(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing id for update", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	req := dto.RequestTeam{}

	// Decode request
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = services.UpdateTeam(uint(id), req)

	// Maybe this should be changed to something else
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Successfull update
	w.WriteHeader(http.StatusNoContent)
}

func DeleteTeam(w http.ResponseWriter, r *http.Request) {

}
