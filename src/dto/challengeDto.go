package dto

import (
	"server/models"
	"time"
)

type ChallengeCreateDto struct {
	Name        string            `json:"name"        validate:"required,min=3"`
	Description string            `json:"description"`
	Sport       string            `json:"sport"       validate:"required,is-valid-sport"`
	Location    LocationCreateDto `json:"location"`
	CreatorId   uint              `json:"creator_id"    validate:"required"`
	IsIndoor    bool              `json:"is_indoor"`
	IsPublic    bool              `json:"is_public"`
	PlayFor     string            `json:"play_for"`
	HasCost     bool              `json:"has_cost"`
	Comment     string            `json:"comment"`
	TeamSize    *int              `json:"team_size"`
	Users       []uint            `json:"users"`
	Teams       []uint            `json:"teams"`
	Date        time.Time         `json:"date"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
}

type ChallengeResponseDto struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Sport       string              `json:"sport"`
	Location    LocationResponseDto `json:"location"`
	Creator     UserResponseDto     `json:"creator"`
	Users       []UserResponseDto   `json:"users"`
	Teams       []TeamResponseDto   `json:"teams"`
	IsIndoor    bool                `json:"is_indoor"`
	IsPublic    bool                `json:"is_public"`
	IsCompleted bool                `json:"is_completed"`
	PlayFor     string              `json:"play_for"`
	HasCost     bool                `json:"has_cost"`
	Comment     string              `json:"comment"`
	TeamSize    *int                `json:"team_size"`
	Date        time.Time           `json:"date"`
	StartTime   time.Time           `json:"time"`
	EndTime     time.Time           `json:"end_time"`
}

func ChallengeCreateDtoToModel(t ChallengeCreateDto) models.Challenge {
	var endTime *time.Time
	if !t.EndTime.IsZero() {
		endTime = &t.EndTime
	}
	return models.Challenge{
		Name:        t.Name,
		Description: t.Description,
		Sport:       t.Sport,
		Location:    LocationCreateDtoToModel(t.Location),
		CreatorID:   t.CreatorId,
		IsIndoor:    t.IsIndoor,
		IsPublic:    t.IsPublic,
		IsCompleted: false,
		PlayFor:     &t.PlayFor,
		HasCost:     t.HasCost,
		Comment:     &t.Comment,
		TeamSize:    t.TeamSize,
		Date:        t.Date,
		StartTime:   t.StartTime,
		EndTime:     endTime,
	}
}

func ToChallengeResponseDto(t models.Challenge) ChallengeResponseDto {
	creator := ToUserResponseDto(t.Creator)
	users := make([]UserResponseDto, len(t.Users))
	for i, user := range t.Users {
		users[i] = ToUserResponseDto(user)
	}
	teams := make([]TeamResponseDto, len(t.Teams))
	for i, team := range t.Teams {
		teams[i] = ToTeamResponseDto(team)
	}
	var endTime time.Time
	if t.EndTime != nil {
		endTime = *t.EndTime
	}
	return ChallengeResponseDto{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		Sport:       t.Sport,
		Location:    ToLocationResponseDto(t.Location),
		Creator:     creator,
		Users:       users,
		Teams:       teams,
		IsIndoor:    t.IsIndoor,
		IsPublic:    t.IsPublic,
		IsCompleted: t.IsCompleted,
		PlayFor:     *t.PlayFor,
		HasCost:     t.HasCost,
		Comment:     *t.Comment,
		TeamSize:    t.TeamSize,
		Date:        t.Date,
		StartTime:   t.StartTime,
		EndTime:     endTime,
	}
}
