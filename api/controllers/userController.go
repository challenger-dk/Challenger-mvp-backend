package controllers

import (
	"encoding/json"
	"net/http"
	"server/api/controllers/helpers"
	"server/common/appError"
	"server/common/dto"
	"server/common/middleware"
	"server/common/models"
	"server/common/services"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Authenticated user
	authUser, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	// Query params
	searchQuery := helpers.GetQueryParamOptional(r, "q")
	limit := helpers.GetQueryInt(r, "limit", 20)
	cursorStr := helpers.GetQueryParamOptional(r, "cursor")

	// Clamp limit (important for protection)
	if limit < 1 {
		limit = 1
	}
	if limit > 50 {
		limit = 50
	}

	// Decode cursor (if provided)
	var cursor *services.UserCursor
	if cursorStr != "" {
		var decoded services.UserCursor
		if err := helpers.DecodeCursor(cursorStr, &decoded); err != nil {
			appError.HandleError(w, appError.ErrBadRequest)
			return
		}
		cursor = &decoded
	}

	// Fetch users (search + pagination)
	users, nextCursor, err := services.GetUsers(
		authUser.ID,
		searchQuery,
		limit,
		cursor,
	)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to DTOs
	out := make([]dto.UserResponseDto, len(users))
	for i, u := range users {
		out[i] = dto.ToUserResponseDto(u)
	}

	// Encode next cursor (if any)
	var nextCursorStr *string
	if nextCursor != nil {
		encoded, err := helpers.EncodeCursor(*nextCursor)
		if err != nil {
			appError.HandleError(w, err)
			return
		}
		nextCursorStr = &encoded
	}

	// Final response
	response := dto.UsersSearchResponse{
		Users:      out,
		NextCursor: nextCursorStr,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)

	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	// Reload user with favorite sports preloaded
	userWithSports, err := services.GetUserByIDWithSettings(user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(dto.ToUserResponseDto(*userWithSports))
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)

	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Use GetVisibleUser to respect blocking rules
	targetUser, err := services.GetVisibleUser(user.ID, id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(dto.ToPublicUserDtoResponse(*targetUser))
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetCurrentUserSettings(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)

	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	settings, err := services.GetUserSettings(user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(*settings)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetInCommonStats(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	targetID, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	stats, err := services.GetInCommonStats(user.ID, targetID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(stats)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetFriends(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)

	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	friends, err := services.GetFriends(user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	response := make([]dto.UserResponseDto, len(friends))
	for i, friend := range friends {
		response[i] = dto.ToUserResponseDto(friend)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func GetSuggestedFriends(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)

	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	suggestedFriends, err := services.GetSuggestedFriends(user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	response := make([]dto.PublicUserDtoResponse, len(suggestedFriends))
	for i, friend := range suggestedFriends {
		response[i] = dto.ToPublicUserDtoResponse(friend)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	req := dto.UserUpdateDto{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.UpdateUser(user.ID, req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().
		Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	req := dto.UserSettingsUpdateDto{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.UpdateUserSettings(user.ID, req)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetParamId(r)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	err = services.DeleteUser(id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func RemoveFriend(w http.ResponseWriter, r *http.Request) {
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

	err = services.RemoveFriend(user.ID, id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func BlockUser(w http.ResponseWriter, r *http.Request) {
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

	err = services.BlockUser(user.ID, id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UnblockUser(w http.ResponseWriter, r *http.Request) {
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

	err = services.UnblockUser(user.ID, id)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
