package controllers

import (
	"github.com/gin-gonic/gin"
	"liango/app/responses"
	"liango/app/services"
	"liango/app/validations"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{service: services.NewAuthService()}
}

// Register creates a new user account.
// POST /auth/register
func (ctrl *AuthController) Register(c *gin.Context) {
	var req validations.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.UnprocessableEntity(c, "Validation failed", err.Error())
		return
	}

	user, err := ctrl.service.Register(req)
	if err != nil {
		if err.Error() == "email already registered" {
			responses.BadRequest(c, err.Error(), nil)
			return
		}
		responses.InternalError(c, "Failed to register user")
		return
	}

	responses.Created(c, "User registered successfully", gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
}

// Login authenticates a user and returns token pair.
// POST /auth/login
func (ctrl *AuthController) Login(c *gin.Context) {
	var req validations.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.UnprocessableEntity(c, "Validation failed", err.Error())
		return
	}

	tokens, user, err := ctrl.service.Login(req)
	if err != nil {
		responses.Unauthorized(c, err.Error())
		return
	}

	responses.Success(c, "Login successful", gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_at":    tokens.ExpiresAt,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// Refresh issues a new token pair using a valid refresh token.
// POST /auth/refresh
func (ctrl *AuthController) Refresh(c *gin.Context) {
	var req validations.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.UnprocessableEntity(c, "Validation failed", err.Error())
		return
	}

	tokens, err := ctrl.service.RefreshToken(req.RefreshToken)
	if err != nil {
		responses.Unauthorized(c, err.Error())
		return
	}

	responses.Success(c, "Token refreshed successfully", gin.H{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
		"expires_at":    tokens.ExpiresAt,
	})
}

// Logout revokes the refresh token.
// POST /auth/logout
func (ctrl *AuthController) Logout(c *gin.Context) {
	var req validations.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.UnprocessableEntity(c, "Validation failed", err.Error())
		return
	}

	_ = ctrl.service.Logout(req.RefreshToken)

	responses.Success(c, "Logged out successfully", nil)
}

// Me returns the currently authenticated user.
// GET /auth/me
func (ctrl *AuthController) Me(c *gin.Context) {
	responses.Success(c, "Authenticated user", gin.H{
		"user_id": c.MustGet("user_id"),
		"email":   c.MustGet("email"),
		"role":    c.MustGet("role"),
	})
}
