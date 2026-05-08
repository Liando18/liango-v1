package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"liango/app/controllers"
	"liango/app/middlewares"
)

// SetupRoutes registers all application routes.
func SetupRoutes(r *gin.Engine) {
	// Root welcome endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to " + os.Getenv("APP_NAME") + " version " + os.Getenv("APP_VERSION"),
		})
	})

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		authCtrl := controllers.NewAuthController()
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authCtrl.Register)
			auth.POST("/login", authCtrl.Login)
			auth.POST("/refresh", authCtrl.Refresh)
			auth.POST("/logout", authCtrl.Logout)

			// Protected auth routes
			meRoutes := auth.Group("/")
			meRoutes.Use(middlewares.JWTMiddleware())
			{
				meRoutes.GET("/me", authCtrl.Me)
			}
		}

		// Resource routes — add more here as your project grows
		RegisterExampleRoutes(v1)
	}
}
