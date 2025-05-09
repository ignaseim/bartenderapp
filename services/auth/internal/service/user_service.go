package service

import (
	"errors"
	"strings"

	"github.com/ignaseim/bartenderapp/services/auth/internal/repository"
	"github.com/ignaseim/bartenderapp/services/pkg/auth"
	"github.com/ignaseim/bartenderapp/services/pkg/models"
)

// UserService handles user-related operations
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}

// GetCurrentUser gets the currently authenticated user
func (s *UserService) GetCurrentUser(claims *auth.Claims) (*models.User, error) {
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""
	return user, nil
}

// List retrieves all users, optionally filtered by role
func (s *UserService) List(role string) ([]models.User, error) {
	users, err := s.userRepo.List(role)
	if err != nil {
		return nil, err
	}

	// Don't return password hashes
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}

// Create adds a new user
func (s *UserService) Create(user *models.User, claims *auth.Claims) error {
	// Validate role
	if !isValidRole(user.Role) {
		return errors.New("invalid role")
	}

	// Only admins can create other admin users
	if user.Role == "admin" && (claims == nil || claims.Role != "admin") {
		return errors.New("only admins can create admin users")
	}

	// Validate email format (basic check)
	if !strings.Contains(user.Email, "@") {
		return errors.New("invalid email format")
	}

	// Hash password
	hashedPassword, err := HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	// Create user
	err = s.userRepo.Create(user)
	if err != nil {
		return err
	}

	// Don't return password hash
	user.PasswordHash = ""
	return nil
}

// Update updates an existing user
func (s *UserService) Update(user *models.User, claims *auth.Claims) error {
	// Validate that the user exists
	existingUser, err := s.userRepo.GetByID(user.ID)
	if err != nil {
		return err
	}

	// Only admins can change roles
	if existingUser.Role != user.Role && (claims == nil || claims.Role != "admin") {
		return errors.New("only admins can change user roles")
	}

	// Only admins can update other admin users
	if existingUser.Role == "admin" && claims.UserID != existingUser.ID && claims.Role != "admin" {
		return errors.New("only admins can update admin users")
	}

	// Check if user is updating themselves
	isSelf := claims.UserID == user.ID

	// If not self and not admin, reject
	if !isSelf && claims.Role != "admin" {
		return errors.New("forbidden: can only update own user or must be admin")
	}

	// If updating password, hash it
	if user.PasswordHash != "" {
		hashedPassword, err := HashPassword(user.PasswordHash)
		if err != nil {
			return err
		}
		user.PasswordHash = hashedPassword
	}

	// Update user
	err = s.userRepo.Update(user)
	if err != nil {
		return err
	}

	// Don't return password hash
	user.PasswordHash = ""
	return nil
}

// Delete removes a user
func (s *UserService) Delete(id int, claims *auth.Claims) error {
	// Get existing user
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check permissions
	if claims.Role != "admin" {
		return errors.New("only admins can delete users")
	}

	// Prevent admin from deleting themselves
	if claims.UserID == id && claims.Role == "admin" {
		return errors.New("admins cannot delete themselves")
	}

	// If trying to delete an admin, require confirmation
	if existingUser.Role == "admin" {
		// In a real application, you might require additional confirmation
		// For now, just allow it
	}

	return s.userRepo.Delete(id)
}

// isValidRole checks if a role is valid
func isValidRole(role string) bool {
	validRoles := []string{"admin", "bartender", "guest"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
} 