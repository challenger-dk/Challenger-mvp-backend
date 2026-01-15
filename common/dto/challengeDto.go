package dto

import (
	"server/common/models"
	"strings"
	"time"
)

// getParticipantsValue returns the participants value from either field
func getParticipantsValue(participants *int, participantCount *int) *int {
	if participants != nil {
		return participants
	}
	return participantCount
}

type ChallengeCreateDto struct {
	Name             string            `json:"name"        validate:"sanitize,required,min=3"`
	Description      string            `json:"description" validate:"sanitize"`
	Sport            string            `json:"sport"       validate:"sanitize,required,is-valid-sport"`
	Location         LocationCreateDto `json:"location"`
	IsIndoor         bool              `json:"is_indoor"`
	IsPublic         bool              `json:"is_public"`
	Status           string            `json:"status" validate:"sanitize"`
	Type             string            `json:"type" validate:"sanitize"`
	PlayFor          string            `json:"play_for"    validate:"sanitize"`
	HasCost          bool              `json:"has_cost"`
	Comment          string            `json:"comment"     validate:"sanitize"`
	TeamSize         *int              `json:"team_size"`
	Distance         *float64          `json:"distance"`
	Participants     *int              `json:"participants"`
	ParticipantCount *int              `json:"participant_count"` // Alternative field name from frontend
	Users            []uint            `json:"users"`
	Teams            []uint            `json:"teams"`
	Date             time.Time         `json:"date"`
	StartTime        time.Time         `json:"start_time"`
	EndTime          time.Time         `json:"end_time"`
}

type ChallengeResponseDto struct {
	ID           uint                    `json:"id"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description"`
	Sport        string                  `json:"sport"`
	Location     LocationResponseDto     `json:"location"`
	Creator      PublicUserDtoResponse   `json:"creator"`
	Users        []PublicUserDtoResponse `json:"users"`
	Teams        []TeamResponseDto       `json:"teams"`
	IsIndoor     bool                    `json:"is_indoor"`
	IsPublic     bool                    `json:"is_public"`
	IsCompleted  bool                    `json:"is_completed"`
	Status       string                  `json:"status" validate:"sanitize"`
	Type         string                  `json:"type" validate:"sanitize"`
	PlayFor      string                  `json:"play_for"`
	HasCost      bool                    `json:"has_cost"`
	Comment      string                  `json:"comment"`
	TeamSize     *int                    `json:"team_size"`
	Distance     *float64                `json:"distance"`
	Participants *int                    `json:"participants"`
	Date         time.Time               `json:"date"`
	StartTime    time.Time               `json:"start_time"`
	EndTime      time.Time               `json:"end_time"`
}

func ChallengeCreateDtoToModel(t ChallengeCreateDto) models.Challenge {
	var endTime *time.Time
	if !t.EndTime.IsZero() {
		endTime = &t.EndTime
	}

	// Set status, defaulting to "pending" if not provided
	status := models.ChallengeStatusPending
	if strings.TrimSpace(t.Status) != "" {
		status = models.ChallengeStatus(strings.TrimSpace(t.Status))
	}

	// Set type, defaulting to "open-for-all" if not provided
	// Normalize the type value (trim whitespace and convert to lowercase for comparison)
	typeValue := strings.TrimSpace(t.Type)
	challengeType := models.ChallengeTypeOpenForAll
	if typeValue != "" {
		// Normalize to match the constant values exactly
		switch strings.ToLower(typeValue) {
		case "run-cycling":
			challengeType = models.ChallengeTypeRunCycling
		case "team-vs-team":
			challengeType = models.ChallengeTypeTeamVsTeam
		case "open-for-all":
			challengeType = models.ChallengeTypeOpenForAll
		default:
			// If it doesn't match, use the provided value as-is (will fail DB constraint if invalid)
			challengeType = models.ChallengeType(typeValue)
		}
	}

	return models.Challenge{
		Name:         t.Name,
		Description:  t.Description,
		Sport:        t.Sport,
		Location:     LocationCreateDtoToModel(t.Location),
		IsIndoor:     t.IsIndoor,
		IsPublic:     t.IsPublic,
		IsCompleted:  false,
		Status:       status,
		Type:         challengeType,
		PlayFor:      &t.PlayFor,
		HasCost:      t.HasCost,
		Comment:      &t.Comment,
		TeamSize:     t.TeamSize,
		Distance:     t.Distance,
		Participants: getParticipantsValue(t.Participants, t.ParticipantCount),
		Date:         t.Date,
		StartTime:    t.StartTime,
		EndTime:      endTime,
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
		ID:           t.ID,
		Name:         t.Name,
		Description:  t.Description,
		Sport:        t.Sport,
		Location:     ToLocationResponseDto(t.Location),
		Creator:      creator,
		Users:        users,
		Teams:        teams,
		IsIndoor:     t.IsIndoor,
		IsPublic:     t.IsPublic,
		IsCompleted:  t.IsCompleted,
		Status:       string(t.Status),
		Type:         string(t.Type),
		PlayFor:      playFor,
		HasCost:      t.HasCost,
		Comment:      comment,
		TeamSize:     t.TeamSize,
		Distance:     t.Distance,
		Participants: t.Participants,
		Date:         t.Date,
		StartTime:    t.StartTime,
		EndTime:      endTime,
	}
}
