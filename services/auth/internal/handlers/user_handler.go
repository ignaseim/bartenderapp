package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ignaseim/bartenderapp/services/auth/internal/service"
	"github.com/ignaseim/bartenderapp/services/pkg/auth"
	"github.com/ignaseim/bartenderapp/services/pkg/middleware"
	"github.com/ignaseim/bartenderapp/services/pkg/models"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser handles requests to get a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get claims from context
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if user is requesting their own data or is admin
	if claims.UserID != id && claims.Role != "admin" {
		middleware.RespondWithError(w, http.StatusForbidden, "Forbidden - Can only view own user or must be admin")
		return
	}

	// Get user
	user, err := h.userService.GetByID(id)
	if err != nil {
		middleware.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, user)
}

// GetCurrentUser handles requests to get the current user
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get current user
	user, err := h.userService.GetCurrentUser(claims)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, user)
}

// ListUsers handles requests to list all users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Only admins can list all users
	if claims.Role != "admin" {
		middleware.RespondWithError(w, http.StatusForbidden, "Forbidden - Requires admin role")
		return
	}

	// Get role filter from query parameters
	role := r.URL.Query().Get("role")

	// Get users
	users, err := h.userService.List(role)
	if err != nil {
		middleware.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, users)
}

// CreateUser handles requests to create a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Only admins can create users
	if claims.Role != "admin" {
		middleware.RespondWithError(w, http.StatusForbidden, "Forbidden - Requires admin role")
		return
	}

	// Parse request body
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create user
	if err := h.userService.Create(&user, claims); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusCreated, user)
}

// UpdateUser handles requests to update an existing user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get claims from context
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse request body
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set user ID from URL
	user.ID = id

	// Update user
	if err := h.userService.Update(&user, claims); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, user)
}

// DeleteUser handles requests to delete a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Get claims from context
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		middleware.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Delete user
	if err := h.userService.Delete(id, claims); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusNoContent, nil)
} 