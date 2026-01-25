package config

import (
	"os"
)

type Config struct {
	AppName string
	AppEnv  string
	AppPort string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// RabbitMQ
	RabbitMQURL string

	// MTN MoMo
	MTNEnvironment     string // "sandbox" or "production"
	MTNBaseURL         string
	MTNSubscriptionKey string
	MTNAPIUser         string
	MTNAPIKey          string
	MTNCallbackURL     string
}

func Load() *Config {
	env := getEnv("MTN_ENVIRONMENT", "sandbox")

	baseURL := "https://sandbox.momodeveloper.mtn.com"
	if env == "production" {
		baseURL = "https://momodeveloper.mtn.com"
	}

	return &Config{
		AppName: getEnv("APP_NAME", "fa-payment-service"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8004"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "farmagent"),
		DBPassword: getEnv("DB_PASSWORD", "farmagent_secret"),
		DBName:     getEnv("DB_NAME", "farmagent_payments"),

		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),

		MTNEnvironment:     env,
		MTNBaseURL:         getEnv("MTN_BASE_URL", baseURL),
		MTNSubscriptionKey: getEnv("MTN_SUBSCRIPTION_KEY", ""),
		MTNAPIUser:         getEnv("MTN_API_USER", ""),
		MTNAPIKey:          getEnv("MTN_API_KEY", ""),
		MTNCallbackURL:     getEnv("MTN_CALLBACK_URL", "http://localhost:8004/api/v1/payments/callback"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
