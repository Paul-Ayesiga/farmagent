package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/farmagent/fa-auth-service/internal/config"
	"github.com/farmagent/fa-auth-service/internal/handlers"
	"github.com/farmagent/fa-auth-service/internal/models"
	"github.com/farmagent/fa-auth-service/internal/repository"
	"github.com/farmagent/fa-auth-service/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByPhone(phone string) (*models.User, error) {
	args := m.Called(phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmailOrPhone(identifier string) (*models.User, error) {
	args := m.Called(identifier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByEmail(email string) bool {
	args := m.Called(email)
	return args.Bool(0)
}

func (m *MockUserRepository) ExistsByPhone(phone string) bool {
	args := m.Called(phone)
	return args.Bool(0)
}

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenID string, expiry time.Duration) error {
	args := m.Called(ctx, userID, tokenID, expiry)
	return args.Error(0)
}

func (m *MockTokenRepository) GetRefreshToken(ctx context.Context, tokenID string) (string, error) {
	args := m.Called(ctx, tokenID)
	return args.String(0), args.Error(1)
}

func (m *MockTokenRepository) DeleteRefreshToken(ctx context.Context, tokenID string) error {
	args := m.Called(ctx, tokenID)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockTokenRepository) StorePasswordResetToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Duration) error {
	args := m.Called(ctx, userID, token, expiry)
	return args.Error(0)
}

func (m *MockTokenRepository) GetPasswordResetToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockTokenRepository) DeletePasswordResetToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepository) StoreEmailVerificationToken(ctx context.Context, userID uuid.UUID, token string, expiry time.Duration) error {
	args := m.Called(ctx, userID, token, expiry)
	return args.Error(0)
}

func (m *MockTokenRepository) GetEmailVerificationToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *MockTokenRepository) DeleteEmailVerificationToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

// Test helpers

func setupTestRouter(authService services.AuthService) *chi.Mux {
	handler := handlers.NewAuthHandler(authService)
	r := chi.NewRouter()

	r.Post("/auth/register", handler.Register)
	r.Post("/auth/login", handler.Login)
	r.Post("/auth/refresh", handler.RefreshToken)
	r.Post("/auth/logout", handler.Logout)

	return r
}

func TestRegister_Success(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)

	cfg := &config.Config{
		JWTSecret:        "test-secret",
		JWTAccessExpiry:  15 * time.Minute,
		JWTRefreshExpiry: 7 * 24 * time.Hour,
	}

	// Set up expectations
	email := "test@example.com"
	mockUserRepo.On("ExistsByEmail", email).Return(false)
	mockUserRepo.On("ExistsByPhone", mock.Anything).Return(false)
	mockUserRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)
	mockTokenRepo.On("StoreRefreshToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	authService := services.NewAuthService(cfg, mockUserRepo, mockTokenRepo, nil)
	router := setupTestRouter(authService)

	body := map[string]interface{}{
		"email":      email,
		"password":   "password123",
		"first_name": "Test",
		"last_name":  "User",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Contains(t, response, "user")
	assert.Contains(t, response, "tokens")
}

func TestRegister_MissingEmailAndPhone(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)

	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	authService := services.NewAuthService(cfg, mockUserRepo, mockTokenRepo, nil)
	router := setupTestRouter(authService)

	body := map[string]interface{}{
		"password":   "password123",
		"first_name": "Test",
		"last_name":  "User",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegister_UserAlreadyExists(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)

	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	email := "existing@example.com"
	mockUserRepo.On("ExistsByEmail", email).Return(true)

	authService := services.NewAuthService(cfg, mockUserRepo, mockTokenRepo, nil)
	router := setupTestRouter(authService)

	body := map[string]interface{}{
		"email":      email,
		"password":   "password123",
		"first_name": "Test",
		"last_name":  "User",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestRegister_InvalidRole(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)

	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	email := "test@example.com"
	mockUserRepo.On("ExistsByEmail", email).Return(false)

	authService := services.NewAuthService(cfg, mockUserRepo, mockTokenRepo, nil)
	router := setupTestRouter(authService)

	body := map[string]interface{}{
		"email":      email,
		"password":   "password123",
		"first_name": "Test",
		"last_name":  "User",
		"role":       "invalid_role",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)

	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	mockUserRepo.On("FindByEmailOrPhone", "nonexistent@example.com").Return(nil, repository.ErrUserNotFound)

	authService := services.NewAuthService(cfg, mockUserRepo, mockTokenRepo, nil)
	router := setupTestRouter(authService)

	body := map[string]interface{}{
		"identifier": "nonexistent@example.com",
		"password":   "password123",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestPasswordValidation(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockTokenRepo := new(MockTokenRepository)

	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	authService := services.NewAuthService(cfg, mockUserRepo, mockTokenRepo, nil)
	router := setupTestRouter(authService)

	// Password too short
	body := map[string]interface{}{
		"email":      "test@example.com",
		"password":   "short",
		"first_name": "Test",
		"last_name":  "User",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var response map[string]string
	json.Unmarshal(rec.Body.Bytes(), &response)
	assert.Contains(t, response["error"], "8 characters")
}
