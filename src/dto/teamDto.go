package dto

import (
	"server/models"
)

/*
type Team struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Users     []User    `gorm:"many2many:user_teams;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Creator   User      `gorm:"foreignKey:CreatorID"`
	CreatorID uint      `gorm:"not null"`
}
*/

type TeamCreateDto struct {
	Name      string `json:"name"        validate:"required,min=3"`
	CreatorId uint   `json:"creator_id"  validate:"required"`
}

type TeamUpdateDto struct {
	Name string `json:"name"        validate:"min=3"`
}

type TeamResponseDto struct {
	ID      uint              `json:"id"`
	Name    string            `json:"name"`
	Creator UserResponseDto   `json:"creator"`
	Users   []UserResponseDto `json:"users"`
}

func TeamCreateDtoToModel(t TeamCreateDto) models.Team {
	return models.Team{
		Name:      t.Name,
		CreatorID: t.CreatorId,
	}
}

func TeamUpdateDtoToModel(t TeamUpdateDto) models.Team {
	return models.Team{
		Name: t.Name,
	}
}

func ToTeamResponseDto(t models.Team) TeamResponseDto {
	var users []UserResponseDto
	for _, u := range t.Users {
		users = append(users, ToUserResponseDto(u))
	}

	return TeamResponseDto{
		ID:      t.ID,
		Name:    t.Name,
		Creator: ToUserResponseDto(t.Creator),
		Users:   users,
	}
}
