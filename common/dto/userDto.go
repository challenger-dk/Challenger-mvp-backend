package dto

import (
	"server/common/models"
	"time"
)

type UserCreateDto struct {
	Email          string   `json:"email"           validate:"required,email"`
	Password       string   `json:"password"        validate:"required,min=8"`
	FirstName      string   `json:"first_name"      validate:"required,min=3"`
	LastName       string   `json:"last_name"`
	ProfilePicture string   `json:"profile_picture,omitempty"`
	Bio            string   `json:"bio,omitempty"`
	Age            uint     `json:"age"             validate:"min=1"`
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
	ID             uint                    `json:"id"`
	Email          string                  `json:"email"`
	FirstName      string                  `json:"first_name"`
	LastName       string                  `json:"last_name"`
	ProfilePicture string                  `json:"profile_picture,omitempty"`
	Bio            string                  `json:"bio,omitempty"`
	Age            uint                    `json:"age"`
	FavoriteSports []SportResponseDto      `json:"favorite_sports,omitempty"`
	Friends        []PublicUserDtoResponse `json:"friends,omitempty"`
	Settings       UserSettingsResponseDto `json:"settings"`
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
}

type UserSettingsResponseDto struct {
	NotifyTeamInvite      bool `json:"notify_team_invite"`
	NotifyFriendReq       bool `json:"notify_friend_req"`
	NotifyChallengeInvite bool `json:"notify_challenge_invite"`
	NotifyChallengeUpdate bool `json:"notify_challenge_update"`
}

type UserSettingsUpdateDto struct {
	NotifyTeamInvite      *bool `json:"notify_team_invite"`
	NotifyFriendReq       *bool `json:"notify_friend_req"`
	NotifyChallengeInvite *bool `json:"notify_challenge_invite"`
	NotifyChallengeUpdate *bool `json:"notify_challenge_update"`
}

// Used for
type PublicUserDtoResponse struct {
	ID             uint               `json:"id"`
	FirstName      string             `json:"first_name"`
	LastName       string             `json:"last_name"`
	ProfilePicture string             `json:"profile_picture,omitempty"`
	Bio            string             `json:"bio,omitempty"`
	Age            uint               `json:"age"`
	FavoriteSports []SportResponseDto `json:"favorite_sports,omitempty"`
}

type Login struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func ToFriendDtoResponse(user models.User) PublicUserDtoResponse {
	favoriteSports := make([]SportResponseDto, len(user.FavoriteSports))
	for i, sport := range user.FavoriteSports {
		favoriteSports[i] = ToSportResponseDto(sport)
	}

	return PublicUserDtoResponse{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		ProfilePicture: user.ProfilePicture,
		Bio:            user.Bio,
		FavoriteSports: favoriteSports,
		Age:            user.Age,
	}
}

func ToUserResponseDto(user models.User) UserResponseDto {
	favoriteSports := make([]SportResponseDto, len(user.FavoriteSports))
	for i, sport := range user.FavoriteSports {
		favoriteSports[i] = ToSportResponseDto(sport)
	}

	friends := make([]PublicUserDtoResponse, len(user.Friends))
	for i, friend := range user.Friends {
		friends[i] = ToFriendDtoResponse(friend)
	}

	var settings UserSettingsResponseDto
	if user.Settings != nil {
		settings = ToUserSettingsResponseDto(*user.Settings)
	} else {
		// Default settings when Settings is nil (all notifications enabled)
		settings = UserSettingsResponseDto{
			NotifyTeamInvite:      true,
			NotifyFriendReq:       true,
			NotifyChallengeInvite: true,
			NotifyChallengeUpdate: true,
		}
	}

	return UserResponseDto{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		ProfilePicture: user.ProfilePicture,
		Bio:            user.Bio,
		Age:            user.Age,
		FavoriteSports: favoriteSports,
		Friends:        friends,
		Settings:       settings,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func ToUserSettingsResponseDto(s models.UserSettings) UserSettingsResponseDto {
	return UserSettingsResponseDto{
		NotifyTeamInvite:      s.NotifyTeamInvite,
		NotifyFriendReq:       s.NotifyFriendReq,
		NotifyChallengeInvite: s.NotifyChallengeInvite,
		NotifyChallengeUpdate: s.NotifyChallengeUpdate,
	}
}

func UserSettingsUpdateDtoToModel(s UserSettingsUpdateDto) models.UserSettings {
	return models.UserSettings{
		NotifyTeamInvite:      *s.NotifyTeamInvite,
		NotifyFriendReq:       *s.NotifyFriendReq,
		NotifyChallengeInvite: *s.NotifyChallengeInvite,
		NotifyChallengeUpdate: *s.NotifyChallengeUpdate,
	}
}
