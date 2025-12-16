package dto

import (
	"server/common/models"
	"time"
)

// --- Request DTOs ---

type CreateDirectConversationDto struct {
	OtherUserID uint `json:"other_user_id" validate:"required,min=1"`
}

type CreateGroupConversationDto struct {
	Title          string `json:"title" validate:"required,min=1,max=255"`
	ParticipantIDs []uint `json:"participant_ids" validate:"required,min=1,dive,min=1"`
}

type SendMessageDto struct {
	Content string `json:"content" validate:"required,min=1,max=2000"`
}

type MarkReadDto struct {
	ReadAt *time.Time `json:"read_at,omitempty"`
}

type SyncTeamMembersDto struct {
	MemberIDs []uint `json:"member_ids" validate:"required,dive,min=1"`
}

// --- Response DTOs ---

type ConversationParticipantDto struct {
	UserID     uint                  `json:"user_id"`
	User       PublicUserDtoResponse `json:"user"`
	JoinedAt   time.Time             `json:"joined_at"`
	LastReadAt *time.Time            `json:"last_read_at,omitempty"`
	LeftAt     *time.Time            `json:"left_at,omitempty"`
}

type ConversationResponseDto struct {
	ID           uint                         `json:"id"`
	Type         models.ConversationType      `json:"type"`
	Title        *string                      `json:"title,omitempty"`
	TeamID       *uint                        `json:"team_id,omitempty"`
	CreatedAt    time.Time                    `json:"created_at"`
	UpdatedAt    time.Time                    `json:"updated_at"`
	Participants []ConversationParticipantDto `json:"participants,omitempty"`
}

type ConversationListItemDto struct {
	ID               uint                    `json:"id"`
	Type             models.ConversationType `json:"type"`
	Title            *string                 `json:"title,omitempty"`
	TeamID           *uint                   `json:"team_id,omitempty"`
	TeamName         *string                 `json:"team_name,omitempty"`
	OtherUser        *PublicUserDtoResponse  `json:"other_user,omitempty"`
	ParticipantCount *int                    `json:"participant_count,omitempty"`
	UnreadCount      int64                   `json:"unread_count"`
	LastMessage      *MessageResponseDto     `json:"last_message,omitempty"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

type MessagesPaginationDto struct {
	Messages []MessageResponseDto `json:"messages"`
	HasMore  bool                 `json:"has_more"`
	Total    int64                `json:"total"`
}

// --- Conversion Functions ---

func ToConversationParticipantDto(p models.ConversationParticipant) ConversationParticipantDto {
	return ConversationParticipantDto{
		UserID:     p.UserID,
		User:       ToPublicUserDtoResponse(p.User),
		JoinedAt:   p.JoinedAt,
		LastReadAt: p.LastReadAt,
		LeftAt:     p.LeftAt,
	}
}

func ToConversationResponseDto(c models.Conversation) ConversationResponseDto {
	dto := ConversationResponseDto{
		ID:        c.ID,
		Type:      c.Type,
		Title:     c.Title,
		TeamID:    c.TeamID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}

	if len(c.Participants) > 0 {
		dto.Participants = make([]ConversationParticipantDto, len(c.Participants))
		for i, p := range c.Participants {
			dto.Participants[i] = ToConversationParticipantDto(p)
		}
	}

	return dto
}

func ToConversationListItemDto(c models.Conversation, unreadCount int64, lastMsg *models.Message, currentUserID uint) ConversationListItemDto {
	dto := ConversationListItemDto{
		ID:          c.ID,
		Type:        c.Type,
		Title:       c.Title,
		TeamID:      c.TeamID,
		UnreadCount: unreadCount,
		UpdatedAt:   c.UpdatedAt,
	}

	// For team conversations, add team name
	if c.Type == models.ConversationTypeTeam && c.Team != nil {
		dto.TeamName = &c.Team.Name
	}

	// For direct conversations, add the other user
	if c.Type == models.ConversationTypeDirect && len(c.Participants) > 0 {
		for _, p := range c.Participants {
			if p.UserID != currentUserID {
				otherUser := ToPublicUserDtoResponse(p.User)
				dto.OtherUser = &otherUser
				break
			}
		}
	}

	// For group and team conversations, add participant count
	if c.Type == models.ConversationTypeGroup || c.Type == models.ConversationTypeTeam {
		count := len(c.Participants)
		dto.ParticipantCount = &count
	}

	if lastMsg != nil {
		msgDto := ToMessageResponseDto(*lastMsg)
		dto.LastMessage = &msgDto
	}

	return dto
}
