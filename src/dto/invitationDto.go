package dto

import (
	"server/models"
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
	ResourceType models.InvitationType   `json:"resource_type"`
	Status       models.InvitationStatus `json:"status"`
}

type InvitationCreateDto struct {
	InviterId    uint                  `json:"inviter_id"`
	InviteeId    uint                  `json:"invitee_id"`
	Note         string                `json:"note"`
	ResourceType models.InvitationType `json:"resource_type"`
	ResourceID   uint                  `json:"resource_id"`
}

func ToInvitationModel(dto InvitationCreateDto) models.Invitation {
	return models.Invitation{
		InviterId:    dto.InviterId,
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
