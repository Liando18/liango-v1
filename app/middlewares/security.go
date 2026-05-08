package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// SecurityHeadersMiddleware adds OWASP-recommended security headers.
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Header("Content-Security-Policy", "default-src 'self'")
		if os.Getenv("APP_ENV") != "development" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		c.Next()
	}
}

// RateLimiterMiddleware limits requests to 100/minute per IP (configurable).
func RateLimiterMiddleware() gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return ginlimiter.NewMiddleware(instance)
}

// APIKeyMiddleware validates X-API-Key header when API_KEY env is set.
func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			c.Next()
			return
		}

		key := c.GetHeader("X-API-Key")
		if key != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPWhitelistMiddleware allows only IPs listed in ALLOWED_IPS env.
// If ALLOWED_IPS is empty, all IPs are allowed.
func IPWhitelistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedIPs := os.Getenv("ALLOWED_IPS")
		if allowedIPs == "" {
			c.Next()
			return
		}

		list := strings.Split(allowedIPs, ",")
		clientIP := c.ClientIP()

		for _, ip := range list {
			if strings.TrimSpace(ip) == clientIP {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Access denied from your IP address",
		})
		c.Abort()
	}
}

// RecoveryMiddleware recovers from panics and returns a 500 response.
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
	})
}

// TimeoutMiddleware cancels requests that exceed the given duration.
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		finished := make(chan struct{})
		go func() {
			c.Next()
			close(finished)
		}()

		select {
		case <-ctx.Done():
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"success": false,
				"message": "Request timeout",
			})
			c.Abort()
		case <-finished:
		}
	}
}
