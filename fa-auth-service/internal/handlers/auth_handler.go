package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/farmagent/fa-auth-service/internal/middleware"
	"github.com/farmagent/fa-auth-service/internal/models"
	"github.com/farmagent/fa-auth-service/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Request/Response types

type RegisterInput struct {
	Phone     *string         `json:"phone,omitempty" example:"+256700123456"`
	Email     *string         `json:"email,omitempty" example:"farmer@example.com"`
	Password  string          `json:"password" example:"password123"`
	FirstName string          `json:"first_name" example:"John"`
	LastName  string          `json:"last_name" example:"Mukasa"`
	Role      models.UserRole `json:"role,omitempty" example:"farmer"`
	District  *string         `json:"district,omitempty" example:"Kampala"`
	Language  string          `json:"language,omitempty" example:"en"`
}

type LoginInput struct {
	Identifier string `json:"identifier" example:"farmer@example.com"`
	Password   string `json:"password" example:"password123"`
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token"`
}

type PasswordInput struct {
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password"`
}

type ForgotPasswordInput struct {
	Identifier string `json:"identifier" example:"farmer@example.com"`
}

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type AssignRoleInput struct {
	Role models.UserRole `json:"role" example:"extension_officer"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with phone or email
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterInput true "Registration details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate
	if input.Phone == nil && input.Email == nil {
		respondError(w, http.StatusBadRequest, "phone or email is required")
		return
	}
	if input.Password == "" {
		respondError(w, http.StatusBadRequest, "password is required")
		return
	}
	if len(input.Password) < 8 {
		respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	tokens, user, err := h.authService.Register(r.Context(), services.RegisterRequest{
		Phone:     input.Phone,
		Email:     input.Email,
		Password:  input.Password,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      input.Role,
		District:  input.District,
		Language:  input.Language,
	})
	if err != nil {
		if err == services.ErrEmailExists {
			respondError(w, http.StatusConflict, "email already registered")
			return
		}
		if err == services.ErrPhoneExists {
			respondError(w, http.StatusConflict, "phone number already registered")
			return
		}
		if err == services.ErrUserExists {
			respondError(w, http.StatusConflict, "user already exists")
			return
		}
		if err == services.ErrInvalidRole {
			respondError(w, http.StatusBadRequest, "invalid role, must be: farmer, extension_officer, buyer, or admin")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

// Login godoc
// @Summary Login to the application
// @Description Authenticate with email/phone and password to get JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, user, err := h.authService.Login(r.Context(), services.LoginRequest{
		Identifier: input.Identifier,
		Password:   input.Password,
	})
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RefreshInput true "Refresh token"
// @Success 200 {object} services.TokenPair
// @Failure 401 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var input RefreshInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tokens, err := h.authService.RefreshToken(r.Context(), input.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	respondJSON(w, http.StatusOK, tokens)
}

// Logout godoc
// @Summary Logout and invalidate refresh token
// @Description Invalidate the refresh token to prevent further use
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RefreshInput true "Refresh token to invalidate"
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var input RefreshInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	h.authService.Logout(r.Context(), input.RefreshToken)
	respondJSON(w, http.StatusOK, map[string]string{"message": "logged out successfully"})
}

// GetProfile godoc
// @Summary Get current user profile
// @Description Get the profile of the authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	user, err := h.authService.GetProfile(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body services.UpdateProfileRequest true "Profile updates"
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /auth/me [put]
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var input services.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.authService.UpdateProfile(r.Context(), userID, input)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// ChangePassword godoc
// @Summary Change password
// @Description Change the authenticated user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body PasswordInput true "Password change request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/change-password [post]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var input PasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.ChangePassword(r.Context(), userID, input.OldPassword, input.NewPassword); err != nil {
		if err == services.ErrInvalidCredentials {
			respondError(w, http.StatusBadRequest, "incorrect current password")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to change password")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "password changed successfully"})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send a password reset email/SMS to the user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body ForgotPasswordInput true "User identifier"
// @Success 200 {object} map[string]string
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var input ForgotPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, _ := h.authService.ForgotPassword(r.Context(), input.Identifier)

	response := map[string]string{
		"message": "if the account exists, a reset link has been sent",
	}
	// In development, include the token for testing
	if token != "" {
		response["token"] = token // Remove in production!
	}

	respondJSON(w, http.StatusOK, response)
}

// ResetPassword godoc
// @Summary Reset password with token
// @Description Reset the user's password using the reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body ResetPasswordInput true "Reset token and new password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input ResetPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.ResetPassword(r.Context(), input.Token, input.NewPassword); err != nil {
		respondError(w, http.StatusBadRequest, "invalid or expired reset token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "password reset successfully"})
}

// AssignRole godoc
// @Summary Assign role to user (Admin only)
// @Description Assign a role to a specific user. Requires admin privileges.
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param input body AssignRoleInput true "Role to assign"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /auth/users/{id}/role [put]
func (h *AuthHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var input AssignRoleInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.AssignRole(r.Context(), userID, input.Role); err != nil {
		if err == services.ErrInvalidRole {
			respondError(w, http.StatusBadRequest, "invalid role")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to assign role")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "role assigned successfully"})
}

// ===== Email Verification =====

type VerifyEmailInput struct {
	Token string `json:"token" example:"123456"`
}

type ResendVerificationInput struct {
	Identifier string `json:"identifier" example:"farmer@example.com"`
}

// SendVerificationEmail godoc
// @Summary Send email verification code
// @Description Send a verification code to the authenticated user's email
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/send-verification [post]
func (h *AuthHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	code, err := h.authService.SendVerificationEmail(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := map[string]string{
		"message": "verification email sent",
	}
	// In development, include the code for testing
	if code != "" {
		response["code"] = code // Remove in production!
	}

	respondJSON(w, http.StatusOK, response)
}

// VerifyEmail godoc
// @Summary Verify email with code
// @Description Verify the user's email using the verification code
// @Tags auth
// @Accept json
// @Produce json
// @Param input body VerifyEmailInput true "Verification code"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var input VerifyEmailInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.authService.VerifyEmail(r.Context(), input.Token); err != nil {
		respondError(w, http.StatusBadRequest, "invalid or expired verification code")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "email verified successfully"})
}

// ResendVerificationEmail godoc
// @Summary Resend email verification
// @Description Resend verification email to the user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body ResendVerificationInput true "User identifier"
// @Success 200 {object} map[string]string
// @Router /auth/resend-verification [post]
func (h *AuthHandler) ResendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	var input ResendVerificationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	h.authService.ResendVerificationEmail(r.Context(), input.Identifier)

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "if the account exists, a verification email has been sent",
	})
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
