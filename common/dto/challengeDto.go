package dto

import (
	"server/common/models"
	"time"
)

type ChallengeCreateDto struct {
	Name        string            `json:"name"        validate:"sanitize,required,min=3"`
	Description string            `json:"description" validate:"sanitize"`
	Sport       string            `json:"sport"       validate:"sanitize,required,is-valid-sport"`
	Location    LocationCreateDto `json:"location"`
	IsIndoor    bool              `json:"is_indoor"`
	IsPublic    bool              `json:"is_public"`
	Status      string            `json:"status"      validate:"sanitize"`
	PlayFor     string            `json:"play_for"    validate:"sanitize"`
	HasCost     bool              `json:"has_cost"`
	Comment     string            `json:"comment"     validate:"sanitize"`
	TeamSize    *int              `json:"team_size"`
	Users       []uint            `json:"users"`
	Teams       []uint            `json:"teams"`
	Date        time.Time         `json:"date"`
	StartTime   time.Time         `json:"start_time"`
	EndTime     time.Time         `json:"end_time"`
}

type ChallengeResponseDto struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Sport       string                  `json:"sport"`
	Location    LocationResponseDto     `json:"location"`
	Creator     PublicUserDtoResponse   `json:"creator"`
	Users       []PublicUserDtoResponse `json:"users"`
	Teams       []TeamResponseDto       `json:"teams"`
	IsIndoor    bool                    `json:"is_indoor"`
	IsPublic    bool                    `json:"is_public"`
	IsCompleted bool                    `json:"is_completed"`
	Status      string                  `json:"status"`
	PlayFor     string                  `json:"play_for"`
	HasCost     bool                    `json:"has_cost"`
	Comment     string                  `json:"comment"`
	TeamSize    *int                    `json:"team_size"`
	Date        time.Time               `json:"date"`
	StartTime   time.Time               `json:"start_time"`
	EndTime     time.Time               `json:"end_time"`
}

func ChallengeCreateDtoToModel(t ChallengeCreateDto) models.Challenge {
	var endTime *time.Time
	if !t.EndTime.IsZero() {
		endTime = &t.EndTime
	}

	// Set status, defaulting to "open" if not provided
	status := models.ChallengeStatusOpen
	if t.Status != "" {
		status = models.ChallengeStatus(t.Status)
	}

	return models.Challenge{
		Name:        t.Name,
		Description: t.Description,
		Sport:       t.Sport,
		Location:    LocationCreateDtoToModel(t.Location),
		IsIndoor:    t.IsIndoor,
		IsPublic:    t.IsPublic,
		IsCompleted: false,
		Status:      status,
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
	creator := ToPublicUserDtoResponse(t.Creator)
	users := make([]PublicUserDtoResponse, len(t.Users))
	for i, user := range t.Users {
		users[i] = ToPublicUserDtoResponse(user)
	}
	teams := make([]TeamResponseDto, len(t.Teams))
	for i, team := range t.Teams {
		teams[i] = ToTeamResponseDto(team)
	}
	var endTime time.Time
	if t.EndTime != nil {
		endTime = *t.EndTime
	}
	var playFor string
	if t.PlayFor != nil {
		playFor = *t.PlayFor
	}
	var comment string
	if t.Comment != nil {
		comment = *t.Comment
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
		Status:      string(t.Status),
		PlayFor:     playFor,
		HasCost:     t.HasCost,
		Comment:     comment,
		TeamSize:    t.TeamSize,
		Date:        t.Date,
		StartTime:   t.StartTime,
		EndTime:     endTime,
	}
}
