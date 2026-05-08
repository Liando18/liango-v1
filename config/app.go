package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Name    string
	Version string
	Port    string
	Env     string
}

var App AppConfig

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("[LianGo] Warning: .env file not found, using system environment variables")
	}

	App = AppConfig{
		Name:    getEnv("APP_NAME", "LianGo"),
		Version: getEnv("APP_VERSION", "1.0.0"),
		Port:    getEnv("APP_PORT", "8080"),
		Env:     getEnv("APP_ENV", "development"),
	}

	log.Printf("[LianGo] Environment loaded: %s v%s (%s)\n", App.Name, App.Version, App.Env)
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
