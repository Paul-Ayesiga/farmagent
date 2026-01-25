package main

import (
	"log"
	"net/http"
	"os"

	"github.com/farmagent/fa-notification-service/internal/config"
	"github.com/farmagent/fa-notification-service/internal/handlers"
	"github.com/farmagent/fa-notification-service/internal/queue"
	"github.com/farmagent/fa-notification-service/internal/services"
	"github.com/farmagent/fa-notification-service/internal/workers"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env in development
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Load config
	cfg := config.Load()

	// Connect to RabbitMQ
	rabbitmq, err := queue.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Initialize services
	emailService := services.NewEmailService(cfg)

	// Start workers
	emailWorker := workers.NewEmailWorker(emailService, rabbitmq)
	if err := emailWorker.Start(); err != nil {
		log.Fatalf("Failed to start email worker: %v", err)
	}

	// Initialize handlers
	notificationHandler := handlers.NewNotificationHandler(rabbitmq)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","service":"fa-notification-service"}`))
	})

	// Notification routes
	r.Route("/api/v1/notifications", func(r chi.Router) {
		r.Post("/email", notificationHandler.SendEmail)
		r.Post("/sms", notificationHandler.SendSMS)
		r.Post("/push", notificationHandler.SendPush)
	})

	// Start server
	port := cfg.AppPort
	if port == "" {
		port = "8003"
	}

	log.Printf("🚀 Notification service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
