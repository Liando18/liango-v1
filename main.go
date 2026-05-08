package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"liango/app/middlewares"
	"liango/app/models"
	"liango/app/routes"
	"liango/app/utils"
	"liango/config"
	"liango/database"
)

func main() {
	// 1. Load environment variables
	config.LoadEnv()

	// 2. Initialize logger
	utils.InitLogger()

	// 3. Set Gin mode
	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 4. Connect to database
	database.Connect()

	// 5. Run migrations
	database.Migrate(
		&models.User{},
		&models.Token{},
		&models.Example{},
		// Add more models here as your project grows
	)

	// 6. Setup Gin engine
	r := gin.New()

	// 7. Register global middlewares
	r.Use(middlewares.RecoveryMiddleware())
	r.Use(middlewares.TimeoutMiddleware(30 * time.Second))
	r.Use(middlewares.SecurityHeadersMiddleware())
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.RateLimiterMiddleware())
	r.Use(utils.GinLogger())

	// 8. Register routes
	routes.SetupRoutes(r)

	// 9. Start server
	addr := ":" + config.App.Port
	log.Printf("[LianGo] Server running at http://localhost%s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[LianGo] Failed to start server: %v", err)
	}
}
