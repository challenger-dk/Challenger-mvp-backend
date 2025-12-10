package dto

import (
	"server/common/models"
)

type CreateEmergencyInfoDto struct {
	Name         string `json:"name"           validate:"sanitize,required,min=2"`
	PhoneNumber  string `json:"phone_number"   validate:"sanitize,required,min=5"` // Todo validate phone number format
	Relationship string `json:"relationship"   validate:"sanitize,required,min=2"` // Todo validate relationship format
}

type EmergencyInfoResponseDto struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	PhoneNumber  string `json:"phone_number"`
	Relationship string `json:"relationship"`
}

func CreateEmergencyInfoDtoToModel(e CreateEmergencyInfoDto) models.EmergencyInfo {
	return models.EmergencyInfo{
		Name:         e.Name,
		PhoneNumber:  e.PhoneNumber,
		Relationship: e.Relationship,
	}
}

func ToEmergencyInfoResponseDto(e models.EmergencyInfo) EmergencyInfoResponseDto {
	return EmergencyInfoResponseDto{
		ID:           e.ID,
		Name:         e.Name,
		PhoneNumber:  e.PhoneNumber,
		Relationship: e.Relationship,
	}
}

func EmergencyInfoResponseDtoToModel(e EmergencyInfoResponseDto) models.EmergencyInfo {
	return models.EmergencyInfo{
		Name:         e.Name,
		PhoneNumber:  e.PhoneNumber,
		Relationship: e.Relationship,
	}
}
