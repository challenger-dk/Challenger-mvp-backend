package dto

import "server/models"

type RequestUser struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type ResponseUser struct {
	ID             uint
	Email          string
	FirstName      string
	LastName       string
	ProfilePicture string
}

func ToResponseUser(u models.User) ResponseUser {
	return ResponseUser{
		ID:             u.ID,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Email:          u.Email,
		ProfilePicture: u.ProfilePicture,
	}
}
