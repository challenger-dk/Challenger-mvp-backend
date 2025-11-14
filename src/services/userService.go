package services

import (
	"errors"

	"server/config"
	"server/dto"
	"server/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists    = errors.New("user with this email already exists")
	ErrInvalidSport  = errors.New("invalid sport name")
	ErrSportNotFound = errors.New("sport not found")
)

func GetUsers() ([]models.User, error) {
	var users []models.User

	err := config.DB.Preload("FavoriteSports").
		Find(&users).
		Error

	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(userID uint) (*models.User, error) {
	var user models.User

	err := config.DB.Preload("FavoriteSports").
		First(&user, userID).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(email, password, firstName, lastName string, favoriteSports []string) (*models.User, error) {
	var existingUser models.User

	err := config.DB.Where("email = ?", email).
		First(&existingUser).
		Error

	if err == nil {
		return nil, ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
	}

	err = config.DB.Create(user).Error
	if err != nil {
		return nil, err
	}

	// Associate favorite sports if provided
	if len(favoriteSports) > 0 {
		if err := associateFavoriteSports(user.ID, favoriteSports); err != nil {
			return nil, err
		}
		// Reload user with favorite sports
		if err := config.DB.Preload("FavoriteSports").First(user, user.ID).Error; err != nil {
			return nil, err
		}
	}

	return user, nil
}

func UpdateUser(userID uint, user dto.UserUpdateDto) error {
	var existingUser models.User
	if err := config.DB.Preload("FavoriteSports").First(&existingUser, userID).Error; err != nil {
		return err
	}

	if user.FirstName != "" {
		existingUser.FirstName = user.FirstName
	}

	if user.LastName != "" {
		existingUser.LastName = user.LastName
	}

	if user.ProfilePicture != "" {
		existingUser.ProfilePicture = user.ProfilePicture
	}

	if user.Bio != "" {
		existingUser.Bio = user.Bio
	}

	if err := config.DB.Save(&existingUser).Error; err != nil {
		return err
	}

	// Update favorite sports if provided
	if user.FavoriteSports != nil {
		if err := associateFavoriteSports(userID, user.FavoriteSports); err != nil {
			return err
		}
	}

	return nil
}

func DeleteUser(userID uint) error {
	return config.DB.Delete(&models.User{}, userID).Error
}
