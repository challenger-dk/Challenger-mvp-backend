package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/common/appError"
	"server/common/dto"
	"server/common/middleware"
	"server/common/models"
	"server/common/services"
	"server/common/validator"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// GetUserFromContext extracts user from context
func GetUserFromContext(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

// CreateDirectConversation creates a new direct conversation
func CreateDirectConversation(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	var req dto.CreateDirectConversationDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appError.HandleError(w, err)
		return
	}

	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	conversation, err := services.CreateDirectConversation(user.ID, req.OtherUserID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToConversationResponseDto(*conversation))
}

// CreateGroupConversation creates a new group conversation
func CreateGroupConversation(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	var req dto.CreateGroupConversationDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appError.HandleError(w, err)
		return
	}

	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	conversation, err := services.CreateGroupConversation(user.ID, req.ParticipantIDs, req.Title)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToConversationResponseDto(*conversation))
}

// ListConversations returns all conversations for the current user
func ListConversations(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	conversations, unreadCounts, lastMessages, err := services.ListConversations(user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to DTOs
	response := make([]dto.ConversationListItemDto, len(conversations))
	for i, conv := range conversations {
		response[i] = dto.ToConversationListItemDto(conv, unreadCounts[i], lastMessages[i], user.ID)
	}

	json.NewEncoder(w).Encode(response)
}

// GetConversation returns a single conversation by ID
func GetConversation(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	conversationID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		appError.HandleError(w, appError.ErrMissingIdParam)
		return
	}

	// Check membership
	isMember, err := services.IsConversationMember(uint(conversationID), user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
	if !isMember {
		appError.HandleError(w, appError.ErrNotConversationMember)
		return
	}

	conversation, err := services.GetConversationByID(uint(conversationID))
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(dto.ToConversationResponseDto(*conversation))
}

// GetConversationMessages returns messages for a conversation with pagination
func GetConversationMessages(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	conversationID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		appError.HandleError(w, appError.ErrMissingIdParam)
		return
	}

	// Parse query parameters
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	var beforeMessageID *uint
	if beforeStr := r.URL.Query().Get("before"); beforeStr != "" {
		if beforeID, err := strconv.ParseUint(beforeStr, 10, 32); err == nil {
			id := uint(beforeID)
			beforeMessageID = &id
		}
	}

	messages, hasMore, total, err := services.GetMessages(uint(conversationID), user.ID, limit, beforeMessageID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Convert to DTOs
	messageDtos := make([]dto.MessageResponseDto, len(messages))
	for i, msg := range messages {
		messageDtos[i] = dto.ToMessageResponseDto(msg)
	}

	response := dto.MessagesPaginationDto{
		Messages: messageDtos,
		HasMore:  hasMore,
		Total:    total,
	}

	json.NewEncoder(w).Encode(response)
}

// SendMessage sends a new message in a conversation
func SendMessage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	conversationID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		appError.HandleError(w, appError.ErrMissingIdParam)
		return
	}

	var req dto.SendMessageDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appError.HandleError(w, err)
		return
	}

	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	message, err := services.SendMessage(uint(conversationID), user.ID, req.Content)
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToMessageResponseDto(*message))
}

// MarkConversationRead marks a conversation as read
func MarkConversationRead(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	conversationID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		appError.HandleError(w, appError.ErrMissingIdParam)
		return
	}

	var req dto.MarkReadDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appError.HandleError(w, err)
		return
	}

	readAt := time.Now()
	if req.ReadAt != nil {
		readAt = *req.ReadAt
	}

	if err := services.MarkConversationRead(uint(conversationID), user.ID, readAt); err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTeamConversation returns the conversation for a team
func GetTeamConversation(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.UserContextKey).(*models.User)
	if !ok {
		appError.HandleError(w, appError.ErrUnauthorized)
		return
	}

	teamID, err := strconv.ParseUint(chi.URLParam(r, "teamId"), 10, 32)
	if err != nil {
		appError.HandleError(w, appError.ErrMissingIdParam)
		return
	}

	// Get or create team conversation
	conversation, err := services.EnsureTeamConversation(uint(teamID))
	if err != nil {
		appError.HandleError(w, err)
		return
	}

	// Check if user is a member of this conversation
	isMember, err := services.IsConversationMember(conversation.ID, user.ID)
	if err != nil {
		appError.HandleError(w, err)
		return
	}
	if !isMember {
		appError.HandleError(w, appError.ErrNotConversationMember)
		return
	}

	json.NewEncoder(w).Encode(dto.ToConversationResponseDto(*conversation))
}

// SyncTeamMembers is an internal endpoint to sync team conversation members
func SyncTeamMembers(w http.ResponseWriter, r *http.Request) {
	teamID, err := strconv.ParseUint(chi.URLParam(r, "teamId"), 10, 32)
	if err != nil {
		appError.HandleError(w, appError.ErrMissingIdParam)
		return
	}

	var req dto.SyncTeamMembersDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appError.HandleError(w, err)
		return
	}

	if err := validator.V.Struct(req); err != nil {
		appError.HandleError(w, err)
		return
	}

	if err := services.SyncTeamConversationMembers(uint(teamID), req.MemberIDs); err != nil {
		appError.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
