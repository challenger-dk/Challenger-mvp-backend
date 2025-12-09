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

func GetCurrentUserInvitations(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	invitations, err := services.GetInvitationsByUserId(user.ID)
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
	if err := validator.V.Struct(invitationDto); err != nil {
		appError.HandleError(w, err)
		return
	}

	invitationModel := dto.ToInvitationModel(invitationDto)

	// Set inviter ID from context
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}
	invitationModel.InviterId = user.ID

	err = services.SendInvitation(&invitationModel)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func AcceptInvitation(w http.ResponseWriter, r *http.Request) {
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

	err = services.AcceptInvitation(id, user.ID)
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

	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	err = services.DeclineInvitation(id, user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}
