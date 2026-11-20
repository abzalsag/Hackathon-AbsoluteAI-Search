package config

import (
	"os"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	AIServiceURL string
}

func LoadConfig() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/ai_search?sslmode=disable"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		AIServiceURL: getEnv("AI_SERVICE_URL", "http://localhost:8000"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
