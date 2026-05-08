package services

import (
	"errors"
	"time"

	"liango/app/helpers"
	"liango/app/models"
	"liango/app/repositories"
	"liango/app/validations"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{userRepo: repositories.NewUserRepository()}
}

// Register creates a new user account.
func (s *AuthService) Register(req validations.RegisterRequest) (*models.User, error) {
	// Check if email already exists
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := helpers.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	role := "user"
	if req.Role != "" {
		role = req.Role
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

// Login validates credentials and returns a token pair.
func (s *AuthService) Login(req validations.LoginRequest) (*helpers.TokenPair, *models.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, nil, errors.New("invalid email or password")
	}

	if !helpers.CheckPassword(user.Password, req.Password) {
		return nil, nil, errors.New("invalid email or password")
	}

	tokens, err := helpers.GenerateTokenPair(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return nil, nil, errors.New("failed to generate token")
	}

	// Store refresh token
	refreshTokenRecord := &models.Token{
		UserID:       user.ID.String(),
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	_ = s.userRepo.SaveRefreshToken(refreshTokenRecord)

	return tokens, user, nil
}

// RefreshToken validates a refresh token and issues new token pair.
func (s *AuthService) RefreshToken(refreshToken string) (*helpers.TokenPair, error) {
	// Parse and validate the refresh token
	claims, err := helpers.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Check if token is stored and not revoked
	storedToken, err := s.userRepo.FindRefreshToken(refreshToken)
	if err != nil || storedToken.Revoked {
		return nil, errors.New("refresh token has been revoked")
	}

	// Revoke the old refresh token (rotation)
	_ = s.userRepo.RevokeRefreshToken(refreshToken)

	// Issue new pair
	tokens, err := helpers.GenerateTokenPair(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		return nil, errors.New("failed to generate new token")
	}

	// Store new refresh token
	newRecord := &models.Token{
		UserID:       claims.UserID,
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	_ = s.userRepo.SaveRefreshToken(newRecord)

	return tokens, nil
}

// Logout revokes the user's refresh token.
func (s *AuthService) Logout(refreshToken string) error {
	return s.userRepo.RevokeRefreshToken(refreshToken)
}
