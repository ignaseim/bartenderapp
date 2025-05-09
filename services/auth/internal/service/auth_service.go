package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yourusername/bartenderapp/services/auth/internal/repository"
	"github.com/yourusername/bartenderapp/services/pkg/auth"
	"github.com/yourusername/bartenderapp/services/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(username, password string) (*models.LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate refresh token (simplified for example)
	refreshToken, err := generateRefreshToken(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create response
	response := &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: models.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return response, nil
}

// RefreshToken generates a new JWT token from a refresh token
func (s *AuthService) RefreshToken(refreshToken string) (*models.LoginResponse, error) {
	// Validate refresh token
	claims, err := validateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Get user from database
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new JWT token
	token, err := auth.GenerateToken(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate new refresh token
	newRefreshToken, err := generateRefreshToken(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create response
	response := &models.LoginResponse{
		Token:        token,
		RefreshToken: newRefreshToken,
		User: models.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	return response, nil
}

// VerifyToken verifies if a JWT token is valid
func (s *AuthService) VerifyToken(tokenString string) (*auth.Claims, error) {
	return auth.ValidateToken(tokenString)
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// generateRefreshToken creates a refresh token for a user
func generateRefreshToken(user models.User) (string, error) {
	// Get secret key from environment
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	// Append "_refresh" to make it different from the access token secret
	refreshSecret := secretKey + "_refresh"

	// Create claims with a longer expiration (e.g., 7 days)
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &auth.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "bartenderapp",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validateRefreshToken validates a refresh token and returns the claims
func validateRefreshToken(tokenString string) (*auth.Claims, error) {
	// Get secret key from environment
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET environment variable is not set")
	}

	// Append "_refresh" to make it different from the access token secret
	refreshSecret := secretKey + "_refresh"

	// Parse the token
	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshSecret), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, errors.New("refresh token expired")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
} 