package dto

import (
	"server/models"
)

type RequestTeam struct {
	Name string
}

type ResponseTeam struct {
	ID    uint
	Name  string
	Users []ResponseUser
}

func RequestTeamToModel(t RequestTeam) models.Team {
	return models.Team{
		Name: t.Name,
	}
}

func ToResponseTeam(t models.Team) ResponseTeam {
	var users []ResponseUser

	for _, u := range t.Users {
		users = append(users, ToResponseUser(u))
	}

	return ResponseTeam{
		ID:    t.ID,
		Name:  t.Name,
		Users: users,
	}
}
