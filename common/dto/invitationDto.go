package dto

import (
	"server/common/models"
)

/*
type Invitation struct {
	ID           uint `gorm:"primaryKey"`
	InviterId    uint `gorm:"not null"`
	Inviter      User `gorm:"foreignKey:InviterId"`
	InviteeId    uint `gorm:"not null"`
	Note         string
	ResourceType InvitationType   `gorm:"type:VARCHAR(20);not null;check:resource_type IN ('team')"`
	ResourceID   uint             `gorm:"not null"`
	Status       InvitationStatus `gorm:"type:VARCHAR(20);not null;default:pending;check:status IN ('pending','accepted','declined')"` // Defualt 'pending'
	CreatedAt    time.Time        `gorm:"autoCreateTime"`
	UpdatedAt    time.Time        `gorm:"autoUpdateTime"`
}
*/

type InvitationResponseDto struct {
	ID           uint                    `json:"id"`
	Inviter      UserResponseDto         `json:"inviter"`
	Note         string                  `json:"note"`
	ResourceType models.ResourceType     `json:"resource_type"`
	Status       models.InvitationStatus `json:"status"`
}

type InvitationCreateDto struct {
	InviteeId    uint                `json:"invitee_id"    validate:"required"`
	Note         string              `json:"note"          validate:"sanitize"`
	ResourceType models.ResourceType `json:"resource_type" validate:"required,oneof=team friend challenge"`
	ResourceID   uint                `json:"resource_id"   validate:"required_if=ResourceType team required_if=ResourceType challenge"` // Friend invitations may not need ResourceID
}

func ToInvitationModel(dto InvitationCreateDto) models.Invitation {
	return models.Invitation{
		InviteeId:    dto.InviteeId,
		Note:         dto.Note,
		ResourceType: dto.ResourceType,
		ResourceID:   dto.ResourceID,
	}
}

func ToInvitationResponse(inv models.Invitation) InvitationResponseDto {
	inviter := ToUserResponseDto(inv.Inviter)
	return InvitationResponseDto{
		ID:           inv.ID,
		Inviter:      inviter,
		Note:         inv.Note,
		ResourceType: inv.ResourceType,
		Status:       inv.Status,
	}
}
