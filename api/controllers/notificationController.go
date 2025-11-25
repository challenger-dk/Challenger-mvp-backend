package controllers

import (
	"encoding/json"
	"net/http"
	"server/api/controllers/helpers"
	"server/api/middleware"
	"server/common/appError"
	"server/common/dto"
	"server/common/models"
	"strconv"

	commonServices "server/common/services"
)

func GetMyNotifications(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	// Parse query parameters
	filters := commonServices.NotificationFilters{}

	// "read" filter (e.g. ?read=false)
	if readStr := r.URL.Query().Get("read"); readStr != "" {
		isRead := readStr == "true"
		filters.IsRead = &isRead
	}

	// "limit" filter
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	// "offset" filter
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	notifs, err := commonServices.GetMyNotifications(user.ID, filters)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	response := make([]dto.NotificationResponseDto, len(notifs))
	for i, n := range notifs {
		response[i] = dto.ToNotificationResponseDto(n)
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func MarkRead(w http.ResponseWriter, r *http.Request) {
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

	err = commonServices.MarkNotificationAsRead(id, user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func MarkAllRead(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	err := commonServices.MarkAllNotificationsAsRead(user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
