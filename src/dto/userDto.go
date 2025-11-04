package dto

import (
	"server/models"
	"time"
)

type UserCreateDto struct {
	Email          string   `json:"email"`
	Password       string   `json:"password"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	ProfilePicture string   `json:"profile_picture,omitempty"`
	Bio            string   `json:"bio,omitempty"`
	FavoriteSports []string `json:"favorite_sports,omitempty"`
}

type UserUpdateDto struct {
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	ProfilePicture string   `json:"profile_picture"`
	Bio            string   `json:"bio,omitempty"`
	FavoriteSports []string `json:"favorite_sports,omitempty"`
}

type SportResponseDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UserResponseDto struct {
	ID             uint               `json:"id"`
	Email          string             `json:"email"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
	ProfilePicture string             `json:"profile_picture,omitempty"`
	Bio            string             `json:"bio,omitempty"`
	FavoriteSports []SportResponseDto `json:"favorite_sports,omitempty"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

func ToSportResponseDto(sport models.Sport) SportResponseDto {
	return SportResponseDto{
		ID:   sport.ID,
		Name: sport.Name,
	}
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func ToUserResponseDto(user models.User) UserResponseDto {
	favoriteSports := make([]SportResponseDto, len(user.FavoriteSports))
	for i, sport := range user.FavoriteSports {
		favoriteSports[i] = ToSportResponseDto(sport)
	}

	return UserResponseDto{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		ProfilePicture: user.ProfilePicture,
		Bio:            user.Bio,
		FavoriteSports: favoriteSports,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}
