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
)

func CreateChat(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	var req dto.CreateChatDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appError.HandleError(w, err)
		return
	}

	if err := validate.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	chat, err := services.CreateChat(user.ID, req.UserIDs, req.Name)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(dto.ToChatResponseDto(*chat))
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetMyChats(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	chats, err := services.GetUserChats(user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	resp := make([]dto.ChatResponseDto, len(chats))
	for i, c := range chats {
		resp[i] = dto.ToChatResponseDto(c)
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetChat(w http.ResponseWriter, r *http.Request) {
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

	chat, err := services.GetChatByID(id, user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(dto.ToChatResponseDto(*chat))
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func AddUserToChat(w http.ResponseWriter, r *http.Request) {
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

	var req dto.AddUserToChatDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.AddUserToChat(id, user.ID, req.UserID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
