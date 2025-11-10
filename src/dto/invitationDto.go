package dto

import (
	"server/models"
)

/*
type Invitation struct {
	ID           uint `gorm:"primaryKey"`
	InviterId    uint `gorm:"not null"`
	InviteeId    uint `gorm:"not null"`
	Note         string
	ResourceType InvitationType   `gorm:"type:VARCHAR(20);not null;check:resource_type IN ('team')"`
	ResourceID   uint             `gorm:"not null"`
	Status       InvitationStatus `gorm:"type:VARCHAR(20);not null;check:status IN ('pending','accepted','declined')"`
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

func ToInvitationResponse(inv models.Invitation, inviter UserResponseDto) InvitationResponseDto {
	return InvitationResponseDto{
		ID:           inv.ID,
		Inviter:      inviter,
		Note:         inv.Note,
		ResourceType: inv.ResourceType,
		Status:       inv.Status,
	}
}
