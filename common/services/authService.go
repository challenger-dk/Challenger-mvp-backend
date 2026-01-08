package services

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"server/common/appError"
	"server/common/config"
	"server/common/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func Login(email, password string) (*models.User, string, error) {
	var user models.User

	err := config.DB.Where("email = ?", email).
		Preload("Settings").
		First(&user).
		Error

	if err != nil {
		return nil, "", appError.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", appError.ErrInvalidCredentials
	}

	token, err := GenerateJWTToken(&user)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func GenerateJWTToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func ValidateJWTToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, appError.ErrInvalidToken
	}

	return claims, nil
}

// generateResetCode generates a random 6-digit code for password reset
func generateResetCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%d", num.Int64())
	}
	return code, nil
}

// RequestPasswordReset generates a reset code, stores it in the database, and sends an email
func RequestPasswordReset(email string) error {
	var user models.User

	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		// Don't reveal if email exists or not for security reasons
		// Return success even if user doesn't exist
		return nil
	}

	// Generate 6-digit reset code
	resetCode, err := generateResetCode()
	if err != nil {
		return fmt.Errorf("failed to generate reset code: %w", err)
	}

	// Set expiration to 1 hour from now
	expiresAt := time.Now().Add(1 * time.Hour)

	// Update user with reset code and expiration
	user.PasswordResetCode = resetCode
	user.PasswordResetCodeExpiresAt = &expiresAt

	err = config.DB.Save(&user).Error
	if err != nil {
		return fmt.Errorf("failed to save reset code: %w", err)
	}

	// Send email with reset code
	err = SendPasswordResetEmail(email, resetCode)
	if err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword validates the reset code and updates the user's password
func ResetPassword(email, resetCode, newPassword string) error {
	var user models.User

	err := config.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return appError.ErrInvalidCredentials
	}

	// Check if reset code exists
	if user.PasswordResetCode == "" {
		return appError.ErrInvalidCredentials
	}

	// Check if reset code matches
	if user.PasswordResetCode != resetCode {
		return appError.ErrInvalidCredentials
	}

	// Check if reset code has expired
	if user.PasswordResetCodeExpiresAt == nil || time.Now().After(*user.PasswordResetCodeExpiresAt) {
		return appError.ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password and clear reset code
	user.Password = string(hashedPassword)
	user.PasswordResetCode = ""
	user.PasswordResetCodeExpiresAt = nil

	err = config.DB.Save(&user).Error
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
