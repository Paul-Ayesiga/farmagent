package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/farmagent/fa-auth-service/internal/config"
	"github.com/farmagent/fa-auth-service/internal/events"
	"github.com/farmagent/fa-auth-service/internal/models"
	"github.com/farmagent/fa-auth-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrEmailExists        = errors.New("email already registered")
	ErrPhoneExists        = errors.New("phone number already registered")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrInvalidRole        = errors.New("invalid role")
)

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*TokenPair, *models.User, error)
	Login(ctx context.Context, req LoginRequest) (*TokenPair, *models.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*models.User, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	ForgotPassword(ctx context.Context, identifier string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
	AssignRole(ctx context.Context, userID uuid.UUID, role models.UserRole) error
	// Email verification
	SendVerificationEmail(ctx context.Context, userID uuid.UUID) (string, error)
	VerifyEmail(ctx context.Context, token string) error
	ResendVerificationEmail(ctx context.Context, identifier string) error
}

type RegisterRequest struct {
	Phone     *string         `json:"phone,omitempty"`
	Email     *string         `json:"email,omitempty"`
	Password  string          `json:"password"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Role      models.UserRole `json:"role,omitempty"`
	District  *string         `json:"district,omitempty"`
	Language  string          `json:"language,omitempty"`
}

type LoginRequest struct {
	Identifier string `json:"identifier"` // email or phone
	Password   string `json:"password"`
}

type UpdateProfileRequest struct {
	FirstName *string  `json:"first_name,omitempty"`
	LastName  *string  `json:"last_name,omitempty"`
	District  *string  `json:"district,omitempty"`
	Language  *string  `json:"language,omitempty"`
	FarmSize  *float64 `json:"farm_size,omitempty"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Claims struct {
	UserID uuid.UUID       `json:"user_id"`
	Role   models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type authService struct {
	cfg       *config.Config
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	publisher *events.Publisher
}

func NewAuthService(cfg *config.Config, userRepo repository.UserRepository, tokenRepo repository.TokenRepository, publisher *events.Publisher) AuthService {
	return &authService{cfg: cfg, userRepo: userRepo, tokenRepo: tokenRepo, publisher: publisher}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*TokenPair, *models.User, error) {
	// Check if user exists - with specific error messages
	if req.Email != nil && s.userRepo.ExistsByEmail(*req.Email) {
		return nil, nil, ErrEmailExists
	}
	if req.Phone != nil && s.userRepo.ExistsByPhone(*req.Phone) {
		return nil, nil, ErrPhoneExists
	}

	// Hash password
	hashedPassword := s.hashPassword(req.Password)

	// Validate and set role (default to farmer)
	role := models.RoleFarmer
	if req.Role != "" {
		if !isValidRole(req.Role) {
			return nil, nil, ErrInvalidRole
		}
		role = req.Role
	}

	// Create user
	user := &models.User{
		Phone:     req.Phone,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		District:  req.District,
		Language:  "en",
		Role:      role,
	}
	if req.Language != "" {
		user.Language = req.Language
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	// Send welcome email if email provided
	if req.Email != nil && s.publisher != nil {
		go s.sendWelcomeEmail(user)
	}

	// Send welcome SMS if phone provided (commented out for now)
	// if req.Phone != nil && s.publisher != nil {
	// 	go s.sendWelcomeSMS(user)
	// }

	// Generate tokens
	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return tokens, user, nil
}

func (s *authService) sendWelcomeEmail(user *models.User) {
	if user.Email == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.publisher.PublishEmail(ctx, events.EmailEvent{
		To:       *user.Email,
		Subject:  "Welcome to FarmAgent! 🌱",
		Template: "welcome",
		Data: map[string]string{
			"name": user.FirstName,
		},
	})
	if err != nil {
		log.Printf("Failed to send welcome email: %v", err)
	}
}

// func (s *authService) sendWelcomeSMS(user *models.User) {
// 	if user.Phone == nil {
// 		return
// 	}
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
//
// 	err := s.publisher.PublishSMS(ctx, events.SMSEvent{
// 		Phone:   *user.Phone,
// 		Message: fmt.Sprintf("Welcome to FarmAgent, %s! Start by scanning your crops for disease detection.", user.FirstName),
// 	})
// 	if err != nil {
// 		log.Printf("Failed to send welcome SMS: %v", err)
// 	}
// }

func (s *authService) Login(ctx context.Context, req LoginRequest) (*TokenPair, *models.User, error) {
	user, err := s.userRepo.FindByEmailOrPhone(req.Identifier)
	if err != nil {
		return nil, nil, ErrInvalidCredentials
	}

	if !s.verifyPassword(req.Password, user.Password) {
		return nil, nil, ErrInvalidCredentials
	}

	tokens, err := s.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return tokens, user, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Parse and validate refresh token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Check if refresh token exists in Redis
	tokenID := claims.ID
	userIDStr, err := s.tokenRepo.GetRefreshToken(ctx, tokenID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, _ := uuid.Parse(userIDStr)
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// Delete old refresh token
	s.tokenRepo.DeleteRefreshToken(ctx, tokenID)

	// Generate new token pair
	return s.generateTokenPair(ctx, user)
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil
	}
	return s.tokenRepo.DeleteRefreshToken(ctx, claims.ID)
}

func (s *authService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *authService) UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.District != nil {
		user.District = req.District
	}
	if req.Language != nil {
		user.Language = *req.Language
	}
	if req.FarmSize != nil {
		user.FarmSize = req.FarmSize
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !s.verifyPassword(oldPassword, user.Password) {
		return ErrInvalidCredentials
	}

	user.Password = s.hashPassword(newPassword)
	return s.userRepo.Update(user)
}

func (s *authService) ForgotPassword(ctx context.Context, identifier string) (string, error) {
	user, err := s.userRepo.FindByEmailOrPhone(identifier)
	if err != nil {
		return "", nil // Don't reveal if user exists
	}

	// Generate 6-digit OTP code for email (easier to type)
	otpCode := s.generateOTPCode(6)

	// Store OTP code (this is what we'll validate)
	err = s.tokenRepo.StorePasswordResetToken(ctx, user.ID, otpCode, 1*time.Hour)
	if err != nil {
		return "", err
	}

	// Send password reset email with the same OTP code
	if user.Email != nil && s.publisher != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			s.publisher.PublishEmail(ctx, events.EmailEvent{
				To:       *user.Email,
				Subject:  "Reset Your FarmAgent Password",
				Template: "password_reset",
				Data: map[string]string{
					"name": user.FirstName,
					"code": otpCode, // Same 6-digit code
				},
			})
		}()
	}

	// Return the same OTP code in response
	return otpCode, nil
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	userIDStr, err := s.tokenRepo.GetPasswordResetToken(ctx, token)
	if err != nil {
		return ErrInvalidToken
	}

	userID, _ := uuid.Parse(userIDStr)
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.Password = s.hashPassword(newPassword)
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	s.tokenRepo.DeletePasswordResetToken(ctx, token)
	return nil
}

func (s *authService) AssignRole(ctx context.Context, userID uuid.UUID, role models.UserRole) error {
	if !isValidRole(role) {
		return ErrInvalidRole
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.Role = role
	return s.userRepo.Update(user)
}

// ===== Email Verification =====

func (s *authService) SendVerificationEmail(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}

	if user.Email == nil {
		return "", errors.New("user has no email")
	}

	if user.IsVerified {
		return "", errors.New("email already verified")
	}

	// Generate 6-digit verification code
	code := s.generateOTPCode(6)
	err = s.tokenRepo.StoreEmailVerificationToken(ctx, user.ID, code, 24*time.Hour)
	if err != nil {
		return "", err
	}

	// Send verification email
	if s.publisher != nil {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			s.publisher.PublishEmail(ctx, events.EmailEvent{
				To:       *user.Email,
				Subject:  "Verify Your FarmAgent Email",
				Template: "email_verification",
				Data: map[string]string{
					"name": user.FirstName,
					"code": code,
				},
			})
		}()
	}

	return code, nil
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	userIDStr, err := s.tokenRepo.GetEmailVerificationToken(ctx, token)
	if err != nil {
		return ErrInvalidToken
	}

	userID, _ := uuid.Parse(userIDStr)
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	user.IsVerified = true
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	s.tokenRepo.DeleteEmailVerificationToken(ctx, token)
	return nil
}

func (s *authService) ResendVerificationEmail(ctx context.Context, identifier string) error {
	user, err := s.userRepo.FindByEmailOrPhone(identifier)
	if err != nil {
		return nil // Don't reveal if user exists
	}

	if user.IsVerified {
		return nil // Silently ignore if already verified
	}

	_, err = s.SendVerificationEmail(ctx, user.ID)
	return err
}

// Helper functions

func isValidRole(role models.UserRole) bool {
	validRoles := []models.UserRole{
		models.RoleFarmer,
		models.RoleExtensionOfficer,
		models.RoleBuyer,
		models.RoleAdmin,
	}
	for _, r := range validRoles {
		if role == r {
			return true
		}
	}
	return false
}

func (s *authService) generateTokenPair(ctx context.Context, user *models.User) (*TokenPair, error) {
	tokenID := uuid.New().String()

	// Access token
	accessClaims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWTAccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshClaims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWTRefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	// Store refresh token in Redis
	err = s.tokenRepo.StoreRefreshToken(ctx, user.ID, tokenID, s.cfg.JWTRefreshExpiry)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    int64(s.cfg.JWTAccessExpiry.Seconds()),
	}, nil
}

func (s *authService) hashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 32)
	return base64.RawStdEncoding.EncodeToString(append(salt, hash...))
}

func (s *authService) verifyPassword(password, encoded string) bool {
	decoded, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil || len(decoded) < 48 {
		return false
	}
	salt := decoded[:16]
	storedHash := decoded[16:]
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 2, 32)
	return string(hash) == string(storedHash)
}

func (s *authService) generateRandomToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// generateOTPCode generates a numeric OTP code of specified length
func (s *authService) generateOTPCode(length int) string {
	const digits = "0123456789"
	b := make([]byte, length)
	rand.Read(b)
	for i := range b {
		b[i] = digits[int(b[i])%len(digits)]
	}
	return string(b)
}
