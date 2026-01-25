package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	// JWT
	JWTSecret string

	// Service URLs
	AuthServiceURL         string
	NotificationServiceURL string
	CropServiceURL         string
	PaymentServiceURL      string
	AIServiceURL           string

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   int // seconds
}

func Load() *Config {
	return &Config{
		AppName: getEnv("APP_NAME", "fa-api-gateway"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8000"),

		JWTSecret: getEnv("JWT_SECRET", "super-secret-key"),

		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "http://localhost:8001"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8003"),
		CropServiceURL:         getEnv("CROP_SERVICE_URL", "http://localhost:8002"),
		PaymentServiceURL:      getEnv("PAYMENT_SERVICE_URL", "http://localhost:8004"),
		AIServiceURL:           getEnv("AI_SERVICE_URL", "http://localhost:8005"),

		RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvInt("RATE_LIMIT_WINDOW", 60),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}
