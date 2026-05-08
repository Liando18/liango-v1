package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// GenerateTokenPair generates an access token + refresh token.
func GenerateTokenPair(userID, email, role string) (*TokenPair, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	// Access token — short lived (15 minutes default)
	accessExpiry := time.Now().Add(15 * time.Minute)
	accessClaims := TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secret)
	if err != nil {
		return nil, err
	}

	// Refresh token — long lived (7 days default)
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := TokenClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExpiry.Unix(),
	}, nil
}

// ParseToken validates and parses a JWT string.
func ParseToken(tokenString string) (*TokenClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// HashPassword hashes a plain password using bcrypt.
func HashPassword(password string) (string, error) {
	// Imported via helpers/password.go — see below
	return bcryptHash(password)
}

// CheckPassword compares a plain password with a bcrypt hash.
func CheckPassword(hashed, plain string) bool {
	return bcryptCheck(hashed, plain)
}
