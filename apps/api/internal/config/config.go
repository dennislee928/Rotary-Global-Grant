package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	// Server settings
	HTTPAddr string
	AppMode  string // dev, staging, prod

	// Database settings
	DatabaseURL string

	// Redis settings
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// JWT settings
	JWTSecret     string
	JWTExpiration time.Duration

	// Rate limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration

	// CAP settings
	CAPSender string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		HTTPAddr:          getEnv("HTTP_ADDR", ":8080"),
		AppMode:           getEnv("APP_MODE", "dev"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://hive:hive@localhost:5432/hive?sslmode=disable"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           getEnvInt("REDIS_DB", 0),
		JWTSecret:         getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		JWTExpiration:     getEnvDuration("JWT_EXPIRATION", 24*time.Hour),
		RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvDuration("RATE_LIMIT_WINDOW", time.Minute),
		CAPSender:         getEnv("CAP_SENDER", "the-hive@example.invalid"),
	}
}

// IsDev returns true if running in development mode
func (c *Config) IsDev() bool {
	return c.AppMode == "dev"
}

// IsProd returns true if running in production mode
func (c *Config) IsProd() bool {
	return c.AppMode == "prod"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
