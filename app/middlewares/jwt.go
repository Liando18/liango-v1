package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"liango/app/helpers"
	"liango/app/responses"
)

// JWTMiddleware validates the Bearer token and sets user context.
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			responses.Unauthorized(c, "Authorization token is required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := helpers.ParseToken(tokenString)
		if err != nil {
			responses.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RoleMiddleware restricts access to specified roles.
// Usage: middlewares.RoleMiddleware("admin", "editor")
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			responses.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, allowed := range roles {
			if role == allowed {
				c.Next()
				return
			}
		}

		responses.Forbidden(c, "You do not have permission to access this resource")
		c.Abort()
	}
}

// PermissionMiddleware restricts access by permission key.
// Usage: middlewares.PermissionMiddleware(constants.PermissionDeleteAny)
func PermissionMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			responses.Unauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		// Import constants inline to avoid circular dep
		role := userRole.(string)
		permMap := map[string][]string{
			"admin":  {"create:any", "read:any", "update:any", "delete:any"},
			"editor": {"create:any", "read:any", "update:own"},
			"user":   {"read:own", "update:own"},
		}

		perms := permMap[role]
		for _, p := range perms {
			if p == permission {
				c.Next()
				return
			}
		}

		responses.Forbidden(c, "Insufficient permissions")
		c.Abort()
	}
}
