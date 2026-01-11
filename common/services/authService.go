package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"server/common/appError"
	"server/common/config"
	"server/common/models"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
	"gorm.io/gorm"
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

	// OAuth users don't have passwords
	if user.AuthProvider != "" {
		return nil, "", appError.ErrInvalidCredentials
	}

	// Regular users must have a password
	if user.Password == nil {
		return nil, "", appError.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
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
	expirationTime := time.Now().Add(time.Duration(config.AppConfig.JWTExpirationHours) * time.Hour)
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

	// Disable prepared statements to avoid cached plan issues after schema changes
	err := config.DB.Session(&gorm.Session{PrepareStmt: false}).
		Where("email = ?", email).First(&user).Error
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

	// Disable prepared statements to avoid cached plan issues after schema changes
	err := config.DB.Session(&gorm.Session{PrepareStmt: false}).
		Where("email = ?", email).First(&user).Error
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
	hashedPasswordStr := string(hashedPassword)
	user.Password = &hashedPasswordStr
	user.PasswordResetCode = ""
	user.PasswordResetCodeExpiresAt = nil

	err = config.DB.Save(&user).Error
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// getFirebaseAuthClient initializes and returns a Firebase Auth client
func getFirebaseAuthClient() (*auth.Client, error) {
	ctx := context.Background()

	// Initialize Firebase app
	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: config.AppConfig.FirebaseProjectID,
	}, option.WithGRPCConnectionPool(10))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Get Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Firebase Auth client: %w", err)
	}

	return authClient, nil
}

// AuthenticateWithGoogle verifies a Google Firebase ID token and returns user with JWT token
func AuthenticateWithGoogle(idToken string) (*models.User, string, error) {
	authClient, err := getFirebaseAuthClient()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get Firebase Auth client: %w", err)
	}

	ctx := context.Background()

	// Verify the Firebase ID token
	token, err := authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, "", appError.ErrInvalidToken
	}

	// Extract user information from the token
	claims := token.Claims

	email, ok := claims["email"].(string)
	if !ok || email == "" {
		return nil, "", fmt.Errorf("email not found in token")
	}

	// Get or create user
	user, err := getOrCreateOAuthUser(email, "google", claims)
	if err != nil {
		return nil, "", err
	}

	// Generate JWT token
	jwtToken, err := GenerateJWTToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, jwtToken, nil
}

// AuthenticateWithApple verifies an Apple Firebase ID token and returns user with JWT token
func AuthenticateWithApple(idToken string, email, firstName, lastName *string) (*models.User, string, error) {
	authClient, err := getFirebaseAuthClient()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get Firebase Auth client: %w", err)
	}

	ctx := context.Background()

	// Verify the Firebase ID token
	token, err := authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, "", appError.ErrInvalidToken
	}

	// Extract user information from the token
	claims := token.Claims

	tokenEmail, ok := claims["email"].(string)
	if !ok || tokenEmail == "" {
		// If email is not in token, use the provided email (Apple only provides email on first sign-in)
		if email == nil || *email == "" {
			return nil, "", fmt.Errorf("email not found in token or request")
		}
		tokenEmail = *email
	}

	// Get or create user
	user, err := getOrCreateOAuthUser(tokenEmail, "apple", claims)
	if err != nil {
		return nil, "", err
	}

	// Update user with provided name information (Apple only provides this on first sign-in)
	if firstName != nil || lastName != nil {
		updateNeeded := false
		if firstName != nil && *firstName != "" && user.FirstName == "" {
			user.FirstName = *firstName
			updateNeeded = true
		}
		if lastName != nil && *lastName != "" && user.LastName == "" {
			user.LastName = *lastName
			updateNeeded = true
		}
		if updateNeeded {
			err = config.DB.Save(user).Error
			if err != nil {
				return nil, "", fmt.Errorf("failed to update user: %w", err)
			}
		}
	}

	// Generate JWT token
	jwtToken, err := GenerateJWTToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, jwtToken, nil
}

// getOrCreateOAuthUser gets an existing OAuth user or creates a new one
func getOrCreateOAuthUser(email, provider string, claims map[string]interface{}) (*models.User, error) {
	var user models.User

	// Try to find existing user
	err := config.DB.Where("email = ?", email).
		Preload("Settings").
		First(&user).
		Error

	if err == nil {
		// User exists
		// If user is OAuth user with different provider, that's an error
		if user.AuthProvider != "" && user.AuthProvider != provider {
			return nil, fmt.Errorf("email already registered with %s", user.AuthProvider)
		}
		// If user exists but doesn't have auth provider set, update it
		if user.AuthProvider == "" {
			user.AuthProvider = provider
			err = config.DB.Save(&user).Error
			if err != nil {
				return nil, fmt.Errorf("failed to update user auth provider: %w", err)
			}
		}
		return &user, nil
	}

	// User doesn't exist, create new OAuth user
	// Extract name from claims if available
	firstName := ""
	lastName := ""

	if name, ok := claims["name"].(string); ok && name != "" {
		// Try to split name into first and last
		// This is a simple approach - Firebase may provide separate fields
		parts := splitName(name)
		if len(parts) > 0 {
			firstName = parts[0]
		}
		if len(parts) > 1 {
			lastName = parts[1]
		}
	}

	// If name not in claims, try given_name and family_name
	if firstName == "" {
		if givenName, ok := claims["given_name"].(string); ok {
			firstName = givenName
		}
	}
	if lastName == "" {
		if familyName, ok := claims["family_name"].(string); ok {
			lastName = familyName
		}
	}

	// If still no first name, use email prefix as fallback
	if firstName == "" {
		firstName = email
	}

	// Create new user
	newUser := models.User{
		Email:        email,
		Password:     nil, // OAuth users don't have passwords
		AuthProvider: provider,
		FirstName:    firstName,
		LastName:     lastName,
		Settings:     &models.UserSettings{},
	}

	err = config.DB.Omit("FavoriteSports").Create(&newUser).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth user: %w", err)
	}

	// Reload with settings
	err = config.DB.Preload("Settings").First(&newUser, newUser.ID).Error
	if err != nil {
		return nil, fmt.Errorf("failed to reload OAuth user: %w", err)
	}

	return &newUser, nil
}

// splitName splits a full name into first and last name
func splitName(name string) []string {
	parts := make([]string, 0)
	current := ""
	for _, char := range name {
		if char == ' ' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
