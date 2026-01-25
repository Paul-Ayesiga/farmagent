package config

import (
	"os"
)

type Config struct {
	// App
	AppName string
	AppEnv  string
	AppPort string

	// MongoDB
	MongoURI string
	MongoDB  string

	// Redis
	RedisHost string
	RedisPort string

	// RabbitMQ
	RabbitMQURL string

	// Email (SMTP)
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
	SMTPFromName string

	// Africa's Talking (SMS)
	ATApiKey   string
	ATUsername string
	ATSenderID string

	// Firebase (Push)
	FCMServerKey string
}

func Load() *Config {
	return &Config{
		// App
		AppName: getEnv("APP_NAME", "fa-notification-service"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8003"),

		// MongoDB
		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGO_DB", "farmagent_notifications"),

		// Redis
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),

		// RabbitMQ
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),

		// Email
		SMTPHost:     getEnv("SMTP_HOST", "localhost"),
		SMTPPort:     getEnv("SMTP_PORT", "1025"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@farmagent.ug"),
		SMTPFromName: getEnv("SMTP_FROM_NAME", "FarmAgent"),

		// Africa's Talking
		ATApiKey:   getEnv("AT_API_KEY", ""),
		ATUsername: getEnv("AT_USERNAME", "sandbox"),
		ATSenderID: getEnv("AT_SENDER_ID", "FarmAgent"),

		// Firebase
		FCMServerKey: getEnv("FCM_SERVER_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
