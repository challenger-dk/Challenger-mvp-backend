package dto

import (
	"server/common/models"
)

/*
type Team struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"not null"`
	Users     []User    `gorm:"many2many:user_teams;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Creator   User      `gorm:"foreignKey:CreatorID"`
	CreatorID uint      `gorm:"not null"`
}
*/

type TeamCreateDto struct {
	Name     string             `json:"name"               validate:"sanitize,required,min=3"`
	Location *LocationCreateDto `json:"location,omitempty"`
}

type TeamUpdateDto struct {
	Name string `json:"name"        validate:"sanitize,min=3"`
}

type TeamResponseDto struct {
	ID       uint                    `json:"id"`
	Name     string                  `json:"name"`
	Creator  PublicUserDtoResponse   `json:"creator"`
	Location LocationResponseDto     `json:"location"`
	Users    []TeamMemberResponseDto `json:"users"`
	Sports   []SportResponseDto      `json:"sports"`
}

type TeamMemberResponseDto struct {
	User PublicUserDtoResponse `json:"user"`
	Role models.TeamRole       `json:"role"`
}

func TeamCreateDtoToModel(t TeamCreateDto) models.Team {
	team := models.Team{
		Name: t.Name,
	}

	if t.Location != nil {
		locationModel := LocationCreateDtoToModel(*t.Location)
		team.Location = &locationModel
	}

	return team
}

func TeamUpdateDtoToModel(t TeamUpdateDto) models.Team {
	return models.Team{
		Name: t.Name,
	}
}

func ToTeamResponseDto(t models.Team) TeamResponseDto {
	var users []TeamMemberResponseDto
	for _, u := range t.Users {
		users = append(users, TeamMemberResponseDto{
			User: ToPublicUserDtoResponse(u.User),
			Role: u.Role,
		})
	}

	var locationDto LocationResponseDto

	if t.Location != nil {
		locationDto = ToLocationResponseDto(*t.Location)
	}

	return TeamResponseDto{
		ID:       t.ID,
		Name:     t.Name,
		Creator:  ToPublicUserDtoResponse(t.Creator),
		Users:    users,
		Location: locationDto,
	}
}
