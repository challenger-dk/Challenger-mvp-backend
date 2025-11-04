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
	if err := config.DB.Preload("FavoriteSports").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := config.DB.Preload("FavoriteSports").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(email, password, firstName, lastName string, favoriteSports []string) (*models.User, error) {
	var existingUser models.User
	if err := config.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
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

	if err := config.DB.Create(user).Error; err != nil {
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

// associateFavoriteSports validates and associates sports with a user
func associateFavoriteSports(userID uint, sportNames []string) error {
	// Get allowed sports
	allowedSports := models.GetAllowedSports()
	allowedSportsMap := make(map[string]bool)
	for _, sport := range allowedSports {
		allowedSportsMap[sport] = true
	}

	// Validate all sport names are allowed
	for _, sportName := range sportNames {
		if !allowedSportsMap[sportName] {
			return ErrInvalidSport
		}
	}

	// Find or create sports in database
	var sports []models.Sport
	for _, sportName := range sportNames {
		var sport models.Sport
		if err := config.DB.Where("name = ?", sportName).FirstOrCreate(&sport, models.Sport{Name: sportName}).Error; err != nil {
			return err
		}
		sports = append(sports, sport)
	}

	// Replace user's favorite sports
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return err
	}

	return config.DB.Model(&user).Association("FavoriteSports").Replace(sports)
}

func DeleteUser(userID uint) error {
	return config.DB.Delete(&models.User{}, userID).Error
}
