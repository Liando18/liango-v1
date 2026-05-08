package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "postgres"
	}

	var dialector gorm.Dialector

	switch driver {
	case "postgres":
		dialector = postgres.Open(buildPostgresDSN())
	case "mysql":
		dialector = mysql.Open(buildMySQLDSN())
	default:
		log.Fatalf("[LianGo] Unsupported DB_DRIVER: %s (use 'postgres' or 'mysql')", driver)
	}

	gormConfig := &gorm.Config{
		Logger: buildLogger(),
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		log.Fatalf("[LianGo] Failed to connect to database: %v", err)
	}

	// Connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("[LianGo] Failed to get sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Printf("[LianGo] Database connected (%s)\n", driver)
}

func buildPostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
}

func buildMySQLDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func buildLogger() logger.Interface {
	level := logger.Warn
	if os.Getenv("APP_ENV") == "development" {
		level = logger.Info
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  level,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
