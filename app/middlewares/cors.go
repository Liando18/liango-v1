package middlewares

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware configures CORS from ALLOWED_ORIGINS env.
func CORSMiddleware() gin.HandlerFunc {
	originsStr := os.Getenv("ALLOWED_ORIGINS")
	origins := []string{"http://localhost:3000", "http://localhost:8080"}
	if originsStr != "" {
		origins = strings.Split(originsStr, ",")
	}

	config := cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(config)
}
