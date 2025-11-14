package dto

import (
	"server/models"
)

/*
type Challenge struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Sport       string
	Location    string
	CreatorID   uint      `gorm:"not null"`
	Creator     User      `gorm:"foreignKey:CreatorID"`
	Teams       []Team    `gorm:"many2many:challenge_teams;"`
	Users       []User    `gorm:"many2many:user_challenges;"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
*/

type ChallengeCreateDto struct {
	Name        string `json:"name"        validate:"required,min=3"`
	Description string `json:"description"`
	Sport       string `json:"sport"       validate:"required,is-valid-sport"`

	Location  string `json:"location"`
	CreatorId uint   `json:"creator_id" validate:"required"`
}

type ChallengeResponseDto struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Sport       string          `json:"sport"`
	Location    string          `json:"location"`
	Creator     UserResponseDto `json:"creator"`
	//Users []UserResponseDto `json:"users"`
}

func ChallengeCreateDtoToModel(t ChallengeCreateDto) models.Challenge {
	return models.Challenge{
		Name:        t.Name,
		Description: t.Description,
		Sport:       t.Sport,
		Location:    t.Location,
		CreatorID:   t.CreatorId,
	}
}

func ToChallengeResponseDto(t models.Challenge) ChallengeResponseDto {
	creator := ToUserResponseDto(t.Creator)

	return ChallengeResponseDto{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Sport:       t.Sport,
		Location:    t.Location,
		Creator:     creator,
	}
}
