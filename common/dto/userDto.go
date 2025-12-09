package dto

import (
	"server/common/models"
)

// How many upcoming challenges to show in user profile
const UserNextChallengesCount uint = 3

type UserCreateDto struct {
	Email          string   `json:"email"           validate:"sanitize,required,email"`
	Password       string   `json:"password"        validate:"sanitize,required,min=8"`
	FirstName      string   `json:"first_name"      validate:"sanitize,required,min=3"`
	LastName       string   `json:"last_name"        validate:"sanitize"`
	ProfilePicture string   `json:"profile_picture,omitempty" validate:"sanitize"`
	Bio            string   `json:"bio,omitempty" validate:"sanitize"`
	Age            uint     `json:"age"             validate:"min=1"`
	FavoriteSports []string `json:"favorite_sports,omitempty"`
}

type UserUpdateDto struct {
	FirstName      string   `json:"first_name"      validate:"sanitize,min=3"`
	LastName       string   `json:"last_name" validate:"sanitize"`
	ProfilePicture string   `json:"profile_picture" validate:"sanitize"`
	Bio            string   `json:"bio,omitempty" validate:"sanitize"`
	FavoriteSports []string `json:"favorite_sports,omitempty"`
}

type UserResponseDto struct {
	ID                  uint                    `json:"id"`
	Email               string                  `json:"email"`
	FirstName           string                  `json:"first_name"`
	LastName            string                  `json:"last_name"`
	ProfilePicture      string                  `json:"profile_picture,omitempty"`
	Bio                 string                  `json:"bio,omitempty"`
	Age                 uint                    `json:"age"`
	FavoriteSports      []SportResponseDto      `json:"favorite_sports,omitempty"`
	Friends             []PublicUserDtoResponse `json:"friends,omitempty"`
	CompletedChallenges uint                    `json:"completed_challenges"`
	NextChallenges      []ChallengeResponseDto  `json:"next_challenges,omitempty"`
	Settings            UserSettingsResponseDto `json:"settings"`
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
	ID                  uint               `json:"id"`
	FirstName           string             `json:"first_name"`
	LastName            string             `json:"last_name"`
	ProfilePicture      string             `json:"profile_picture,omitempty"`
	Bio                 string             `json:"bio,omitempty"`
	Age                 uint               `json:"age"`
	FavoriteSports      []SportResponseDto `json:"favorite_sports,omitempty"`
	FriendsCount        uint               `json:"friends_count,omitempty"`
	TeamsCount          uint               `json:"teams_count,omitempty"`
	CompletedChallenges uint               `json:"completed_challenges,omitempty"`
}

type Login struct {
	Email    string `json:"email"    validate:"sanitize,required,email"`
	Password string `json:"password" validate:"sanitize,required,min=8"`
}

type CommonStatsDto struct {
	CommonFriendsCount int64      `json:"common_friends_count"`
	CommonTeamsCount   int64      `json:"common_teams_count"`
	CommonSports       []SportDto `json:"common_sports"`
}

func ToPublicUserDtoResponse(user models.User) PublicUserDtoResponse {
	favoriteSports := make([]SportResponseDto, len(user.FavoriteSports))
	for i, sport := range user.FavoriteSports {
		favoriteSports[i] = ToSportResponseDto(sport)
	}

	friendsCount := uint(len(user.Friends))
	teamsCount := uint(len(user.Teams))

	// Count completed challenges
	var completedChallengesCount uint
	for _, challenge := range user.JoinedChallenges {
		if challenge.IsCompleted {
			completedChallengesCount++
		}
	}

	return PublicUserDtoResponse{
		ID:                  user.ID,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		ProfilePicture:      user.ProfilePicture,
		Bio:                 user.Bio,
		FavoriteSports:      favoriteSports,
		Age:                 user.Age,
		FriendsCount:        friendsCount,
		TeamsCount:          teamsCount,
		CompletedChallenges: completedChallengesCount,
	}
}

func ToUserResponseDto(user models.User) UserResponseDto {
	favoriteSports := make([]SportResponseDto, len(user.FavoriteSports))
	for i, sport := range user.FavoriteSports {
		favoriteSports[i] = ToSportResponseDto(sport)
	}

	friends := make([]PublicUserDtoResponse, len(user.Friends))
	for i, friend := range user.Friends {
		friends[i] = ToPublicUserDtoResponse(friend)
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

	// Count completed challenges
	var completedChallengesCount uint
	for _, challenge := range user.JoinedChallenges {
		if challenge.IsCompleted {
			completedChallengesCount++
		}
	}

	// Get next upcoming challenges (limited to UserNextChallengesCount)
	nextChallenges := make([]ChallengeResponseDto, 0, UserNextChallengesCount)
	for _, ch := range user.JoinedChallenges {
		if ch.IsCompleted {
			continue
		}
		if len(nextChallenges) >= int(UserNextChallengesCount) {
			break
		}
		nextChallenges = append(nextChallenges, ToChallengeResponseDto(ch))
	}

	return UserResponseDto{
		ID:                  user.ID,
		Email:               user.Email,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		ProfilePicture:      user.ProfilePicture,
		Bio:                 user.Bio,
		Age:                 user.Age,
		FavoriteSports:      favoriteSports,
		Friends:             friends,
		Settings:            settings,
		CompletedChallenges: completedChallengesCount,
		NextChallenges:      nextChallenges,
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

func UserCreateDtoToModel(u UserCreateDto) models.User {
	favoriteSports := make([]models.Sport, len(u.FavoriteSports))
	for i, sportName := range u.FavoriteSports {
		favoriteSports[i] = models.Sport{Name: sportName}
	}

	return models.User{
		Email:          u.Email,
		Password:       u.Password,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		ProfilePicture: u.ProfilePicture,
		Bio:            u.Bio,
		Age:            u.Age,
		FavoriteSports: favoriteSports,
		Settings:       &models.UserSettings{},
	}
}
