package routes

import (
	"liango/app/constants"
	"liango/app/controllers"
	"liango/app/middlewares"

	"github.com/gin-gonic/gin"
)

// RegisterExampleRoutes registers all routes for the Example resource.
// Copy & rename this for other entities.
func RegisterExampleRoutes(rg *gin.RouterGroup) {
	ctrl := controllers.NewExampleController()

	examples := rg.Group("/examples")
	examples.Use(middlewares.JWTMiddleware())
	{
		// Public (any authenticated user)
		examples.GET("", ctrl.Index)
		examples.GET("/:id", ctrl.Show)

		// Protected (admin/editor only)
		examples.POST("", middlewares.PermissionMiddleware(constants.PermissionCreateAny), ctrl.Store)
		examples.PUT("/:id", middlewares.PermissionMiddleware(constants.PermissionUpdateAny), ctrl.Update)
		examples.DELETE("/:id", middlewares.RoleMiddleware(constants.RoleAdmin), ctrl.Destroy)
	}
}
