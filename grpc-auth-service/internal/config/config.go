package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration properties for the application
type Config struct {
	Port          string
	JWTSecret     string
	JWTExpiration time.Duration
	Environment   string
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() *Config {
	// Attempt to load .env file if available
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️  No .env file found; using environment variables and defaults")
	}

	port := getEnv("PORT", "50051")
	jwtSecret := getEnv("JWT_SECRET", "super-secret-grpc-jwt-key-change-in-production")
	env := getEnv("ENV", "development")

	expStr := getEnv("JWT_EXPIRATION", "24h")
	jwtExp, err := time.ParseDuration(expStr)
	if err != nil {
		jwtExp = 24 * time.Hour
	}

	return &Config{
		Port:          port,
		JWTSecret:     jwtSecret,
		JWTExpiration: jwtExp,
		Environment:   env,
	}
}

// getEnv retrieves an environment variable or returns a fallback value if not set
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
