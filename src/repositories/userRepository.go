package repositories

import (
	"server/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(db *gorm.DB, user *models.User) error
	GetUserByID(db *gorm.DB, id uint) (*models.User, error)
	GetUserByEmail(db *gorm.DB, email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}
