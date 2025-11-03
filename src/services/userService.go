package services

import (
	"server/models"
	"server/repositories"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(db *gorm.DB, id uint) (*models.User, error)
	GetUserByEmail(db *gorm.DB, email string) (*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}