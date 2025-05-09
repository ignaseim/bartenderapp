package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yourusername/bartenderapp/services/auth/internal/service"
	"github.com/yourusername/bartenderapp/services/pkg/middleware"
	"github.com/yourusername/bartenderapp/services/pkg/models"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login requests
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Call service to authenticate user
	resp, err := h.authService.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, resp)
}

// RefreshToken handles token refresh requests
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Call service to refresh token
	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	middleware.RespondWithJSON(w, http.StatusOK, resp)
}

// VerifyToken handles token verification requests
func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Call service to verify token
	claims, err := h.authService.VerifyToken(req.Token)
	if err != nil {
		middleware.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Format expiration time
	var expiresAt time.Time
	if expClaim, err := claims.GetExpirationTime(); err == nil && expClaim != nil {
		expiresAt = expClaim.Time
	}

	// Create response
	resp := map[string]interface{}{
		"valid":      true,
		"user_id":    claims.UserID,
		"username":   claims.Username,
		"role":       claims.Role,
		"expires_at": expiresAt,
	}

	middleware.RespondWithJSON(w, http.StatusOK, resp)
} 