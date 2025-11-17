package controllers

import (
	"encoding/json"
	"net/http"
	"server/appError"
	"server/controllers/helpers"
	"server/dto"
	"server/services"
)

func GetInvitationsByUserId(w http.ResponseWriter, r *http.Request) {
	user_id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	invitations, err := services.GetInvitationsByUserId(user_id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to response DTOs
	response := make([]dto.InvitationResponseDto, len(invitations))
	for i, inv := range invitations {
		response[i] = dto.ToInvitationResponse(inv)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func SendInvitation(w http.ResponseWriter, r *http.Request) {
	var invitationDto dto.InvitationCreateDto

	err := json.NewDecoder(r.Body).Decode(&invitationDto)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Validate
	if err := validate.Struct(invitationDto); err != nil {
		appError.HandleError(w, err)
		return
	}

	invitationModel := dto.ToInvitationModel(invitationDto)
	err = services.SendInvitation(&invitationModel)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

// TODO: Implement a way to tell which user is accepting/declining the invitation
// to avoid unauthorized actions.
func AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.AcceptInvitation(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func DeclineInvitation(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.DeclineInvitation(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}
