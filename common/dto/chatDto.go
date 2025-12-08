package dto

import "server/common/models"

type CreateChatDto struct {
	Name    string `json:"name"`
	UserIDs []uint `json:"user_ids" validate:"required,min=1"`
}

type AddUserToChatDto struct {
	UserID uint `json:"user_id" validate:"required"`
}

type ChatResponseDto struct {
	ID    uint              `json:"id"`
	Name  string            `json:"name"`
	Users []UserResponseDto `json:"users"`
}

func ToChatResponseDto(c models.Chat) ChatResponseDto {
	users := make([]UserResponseDto, len(c.Users))
	for i, u := range c.Users {
		users[i] = ToUserResponseDto(u)
	}

	return ChatResponseDto{
		ID:    c.ID,
		Name:  c.Name,
		Users: users,
	}
}
