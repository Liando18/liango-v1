package database

import (
	"log"

	"gorm.io/gorm"
)

// Migrate runs GORM AutoMigrate for all given models.
// Usage: database.Migrate(&models.User{}, &models.Example{})
func Migrate(models ...interface{}) {
	if err := DB.AutoMigrate(models...); err != nil {
		log.Fatalf("[LianGo] Migration failed: %v", err)
	}
	log.Println("[LianGo] Database migration completed")
}

// GetDB returns the active database instance.
func GetDB() *gorm.DB {
	return DB
}
