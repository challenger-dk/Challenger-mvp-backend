package dto

import (
	"server/models"
	"time"
)

type UserCreateDto struct {
	Email          string   `json:"email"           validate:"required,email"`
	Password       string   `json:"password"        validate:"required,min=8"`
	FirstName      string   `json:"first_name"      validate:"required,min=3"`
	LastName       string   `json:"last_name"`
	ProfilePicture string   `json:"profile_picture,omitempty"`
	Bio            string   `json:"bio,omitempty"`
	FavoriteSports []string `json:"favorite_sports,omitempty"`
}

type UserUpdateDto struct {
	FirstName      string   `json:"first_name"      validate:"min=3"`
	LastName       string   `json:"last_name"`
	ProfilePicture string   `json:"profile_picture"`
	Bio            string   `json:"bio,omitempty"`
	FavoriteSports []string `json:"favorite_sports,omitempty"`
}

type UserResponseDto struct {
	ID             uint       `json:"id"`
	Email          string     `json:"email"`
	FirstName      string     `json:"first_name"`
	LastName       string     `json:"last_name"`
	ProfilePicture string     `json:"profile_picture,omitempty"`
	Bio            string     `json:"bio,omitempty"`
	FavoriteSports []SportDto `json:"favorite_sports,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Login struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func ToUserResponseDto(user models.User) UserResponseDto {
	favoriteSports := make([]SportDto, len(user.FavoriteSports))
	for i, sport := range user.FavoriteSports {
		favoriteSports[i] = ToSportDto(sport)
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
