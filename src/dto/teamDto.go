package dto

import (
	"server/models"
)

type TeamCreateDto struct {
	Name string
}

type TeamResponseDto struct {
	ID    uint
	Name  string
	Users []UserResponseDto
}

func TeamCreateDtoToModel(t TeamCreateDto) models.Team {
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
		ID:    t.ID,
		Name:  t.Name,
		Users: users,
	}
}
