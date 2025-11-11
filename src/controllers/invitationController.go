package controllers

import (
	"encoding/json"
	"net/http"
	"server/controllers/helpers"
	"server/dto"
	"server/services"
)

func GetInvitationsByUserId(w http.ResponseWriter, r *http.Request) {
	user_id := helpers.GetParamId(w, r)
	if user_id == 0 {
		return
	}

	invitations, err := services.GetInvitationsByUserId(user_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response DTOs
	response := make([]dto.InvitationResponseDto, len(invitations))
	for i, inv := range invitations {
		response[i] = dto.ToInvitationResponse(inv)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SendInvitation(w http.ResponseWriter, r *http.Request) {
	var invitationDto dto.InvitationCreateDto

	err := json.NewDecoder(r.Body).Decode(&invitationDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	invitationModel := dto.ToInvitationModel(invitationDto)
	err = services.SendInvitation(&invitationModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamId(w, r)
	if id == 0 {
		return
	}

	err := services.AcceptInvitation(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeclineInvitation(w http.ResponseWriter, r *http.Request) {
	id := helpers.GetParamId(w, r)
	if id == 0 {
		return
	}

	err := services.DeclineInvitation(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
