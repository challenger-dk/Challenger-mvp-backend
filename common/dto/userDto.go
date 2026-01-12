package dto

import (
	"server/common/models"
	"time"
)

// How many upcoming challenges to show in user profile
const UserNextChallengesCount uint = 3

type UserCreateDto struct {
	Email          string    `json:"email"           validate:"sanitize,required,email"`
	Password       string    `json:"password"        validate:"sanitize,required,min=8"`
	FirstName      string    `json:"first_name"      validate:"sanitize,required,min=3"`
	LastName       string    `json:"last_name"        validate:"sanitize"`
	ProfilePicture string    `json:"profile_picture,omitempty" validate:"sanitize"`
	Bio            string    `json:"bio,omitempty" validate:"sanitize"`
	BirthDate      time.Time `json:"birth_date"      validate:"required"`
	City           string    `json:"city"            validate:"sanitize"`
	FavoriteSports []string  `json:"favorite_sports,omitempty"`
}

type UserUpdateDto struct {
	FirstName      string    `json:"first_name"      validate:"sanitize,min=3"`
	LastName       string    `json:"last_name" validate:"sanitize"`
	ProfilePicture string    `json:"profile_picture" validate:"sanitize"`
	Bio            string    `json:"bio,omitempty" validate:"sanitize"`
	BirthDate      time.Time `json:"birth_date"      validate:"required"`
	City           string    `json:"city"            validate:"sanitize"`
	FavoriteSports []string  `json:"favorite_sports,omitempty"`
}

type UserResponseDto struct {
	ID                  uint                       `json:"id"`
	Email               string                     `json:"email"`
	FirstName           string                     `json:"first_name"`
	LastName            string                     `json:"last_name"`
	ProfilePicture      string                     `json:"profile_picture,omitempty"`
	Bio                 string                     `json:"bio,omitempty"`
	BirthDate           time.Time                  `json:"birth_date"`
	City                string                     `json:"city"`
	FavoriteSports      []SportResponseDto         `json:"favorite_sports,omitempty"`
	Friends             []PublicUserDtoResponse    `json:"friends,omitempty"`
	CompletedChallenges uint                       `json:"completed_challenges"`
	NextChallenges      []ChallengeResponseDto     `json:"next_challenges,omitempty"`
	Settings            UserSettingsResponseDto    `json:"settings"`
	EmergencyContacts   []EmergencyInfoResponseDto `json:"emergency_contacts,omitempty"`
	Teams               []TeamResponseDto          `json:"teams,omitempty"`
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

type UsersSearchResponse struct {
	Users      []UserResponseDto `json:"users"`
	NextCursor *string           `json:"next_cursor"`
}

// Used for anyone but the current user
type PublicUserDtoResponse struct {
	ID                  uint                   `json:"id"`
	FirstName           string                 `json:"first_name"`
	LastName            string                 `json:"last_name"`
	ProfilePicture      string                 `json:"profile_picture,omitempty"`
	Bio                 string                 `json:"bio,omitempty"`
	BirthDate           time.Time              `json:"birth_date"`
	City                string                 `json:"city"`
	FavoriteSports      []SportResponseDto     `json:"favorite_sports,omitempty"`
	FriendsCount        uint                   `json:"friends_count,omitempty"`
	TeamsCount          uint                   `json:"teams_count,omitempty"`
	CompletedChallenges uint                   `json:"completed_challenges,omitempty"`
	NextChallenges      []ChallengeResponseDto `json:"next_challenges,omitempty"`
}

type Login struct {
	Email    string `json:"email"    validate:"sanitize,required,email"`
	Password string `json:"password" validate:"sanitize,required,min=8"`
}

type RequestPasswordResetDto struct {
	Email string `json:"email" validate:"sanitize,required,email"`
}

type ResetPasswordDto struct {
	Email       string `json:"email"        validate:"sanitize,required,email"`
	ResetCode   string `json:"reset_code"   validate:"sanitize,required,len=6"`
	NewPassword string `json:"new_password" validate:"sanitize,required,min=8"`
}

type GoogleAuthDto struct {
	IDToken string `json:"idToken" validate:"sanitize,required"`
}

type AppleAuthDto struct {
	IDToken   string  `json:"idToken" validate:"sanitize,required"`
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
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

	return PublicUserDtoResponse{
		ID:                  user.ID,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		ProfilePicture:      user.ProfilePicture,
		Bio:                 user.Bio,
		FavoriteSports:      favoriteSports,
		BirthDate:           user.BirthDate,
		City:                user.City,
		FriendsCount:        friendsCount,
		TeamsCount:          teamsCount,
		CompletedChallenges: completedChallengesCount,
		NextChallenges:      nextChallenges,
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

	emergencyContacts := make([]EmergencyInfoResponseDto, len(user.EmergencyContacts))
	for i, contact := range user.EmergencyContacts {
		emergencyContacts[i] = ToEmergencyInfoResponseDto(contact)
	}

	teams := make([]TeamResponseDto, len(user.Teams))
	for i, team := range user.Teams {
		teams[i] = ToTeamResponseDto(team)
	}

	return UserResponseDto{
		ID:                  user.ID,
		Email:               user.Email,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		ProfilePicture:      user.ProfilePicture,
		Bio:                 user.Bio,
		BirthDate:           user.BirthDate,
		City:                user.City,
		FavoriteSports:      favoriteSports,
		Friends:             friends,
		Settings:            settings,
		CompletedChallenges: completedChallengesCount,
		NextChallenges:      nextChallenges,
		EmergencyContacts:   emergencyContacts,
		Teams:               teams,
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

	var passwordPtr *string
	if u.Password != "" {
		passwordPtr = &u.Password
	}

	return models.User{
		Email:          u.Email,
		Password:       passwordPtr,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		ProfilePicture: u.ProfilePicture,
		Bio:            u.Bio,
		BirthDate:      u.BirthDate,
		City:           u.City,
		FavoriteSports: favoriteSports,
		Settings:       &models.UserSettings{},
	}
}
