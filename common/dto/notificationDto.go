package dto

import (
	"server/common/models"
	"time"
)

type NotificationResponseDto struct {
	ID           uint                    `json:"id"`
	Title        string                  `json:"title"`
	Content      string                  `json:"content"`
	Type         models.NotificationType `json:"type"`
	IsRead       bool                    `json:"is_read"`
	IsRelevant   bool                    `json:"is_relevant"`
	Actor        *PublicUserDtoResponse  `json:"actor,omitempty"`
	ResourceID   *uint                   `json:"resource_id,omitempty"`
	ResourceType *string                 `json:"resource_type,omitempty"`
	InvitationID *uint                   `json:"invitation_id,omitempty"`
	CreatedAt    time.Time               `json:"created_at"`
}

func ToNotificationResponseDto(n models.Notification) NotificationResponseDto {
	dto := NotificationResponseDto{
		ID:           n.ID,
		Title:        n.Title,
		Content:      n.Content,
		Type:         n.Type,
		IsRead:       n.IsRead,
		IsRelevant:   n.IsRelevant,
		ResourceID:   n.ResourceID,
		ResourceType: n.ResourceType,
		InvitationID: n.InvitationID,
		CreatedAt:    n.CreatedAt,
	}

	if n.Actor != nil {
		actorDto := ToFriendDtoResponse(*n.Actor)
		dto.Actor = &actorDto
	}

	return dto
}
