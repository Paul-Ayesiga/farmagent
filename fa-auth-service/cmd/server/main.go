package main

import (
	"log"
	"net/http"
	"os"

	"github.com/farmagent/fa-auth-service/internal/config"
	"github.com/farmagent/fa-auth-service/internal/database"
	"github.com/farmagent/fa-auth-service/internal/events"
	"github.com/farmagent/fa-auth-service/internal/handlers"
	"github.com/farmagent/fa-auth-service/internal/middleware"
	"github.com/farmagent/fa-auth-service/internal/models"
	"github.com/farmagent/fa-auth-service/internal/repository"
	"github.com/farmagent/fa-auth-service/internal/services"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/farmagent/fa-auth-service/docs"
)

// @title FarmAgent Auth Service API
// @version 1.0
// @description Authentication service for FarmAgent - handles user registration, login, JWT tokens, and role management.
// @termsOfService http://swagger.io/terms/

// @contact.name FarmAgent Support
// @contact.email support@farmagent.ug

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8001
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format: Bearer {token}

func main() {
	// Load .env in development
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Load config
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Connect to Redis
	redisClient, err := database.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Connect to RabbitMQ
	var publisher *events.Publisher
	publisher, err = events.NewPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("⚠️ Failed to connect to RabbitMQ: %v (emails will not be sent)", err)
		publisher = nil
	} else {
		defer publisher.Close()
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(redisClient)

	// Initialize services
	authService := services.NewAuthService(cfg, userRepo, tokenRepo, publisher)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","service":"fa-auth-service"}`))
	})

	// Auth routes
	r.Route("/api/v1/auth", func(r chi.Router) {
		// Public routes
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.RefreshToken)
		r.Post("/logout", authHandler.Logout)
		r.Post("/forgot-password", authHandler.ForgotPassword)
		r.Post("/reset-password", authHandler.ResetPassword)
		r.Post("/verify-email", authHandler.VerifyEmail)
		r.Post("/resend-verification", authHandler.ResendVerificationEmail)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))
			r.Get("/me", authHandler.GetProfile)
			r.Put("/me", authHandler.UpdateProfile)
			r.Post("/change-password", authHandler.ChangePassword)
			r.Post("/send-verification", authHandler.SendVerificationEmail)
		})

		// Admin routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))
			r.Use(middleware.RequireRole(models.RoleAdmin))
			r.Put("/users/{id}/role", authHandler.AssignRole)
		})
	})

	// Start server
	port := cfg.AppPort
	if port == "" {
		port = "8001"
	}

	log.Printf("🚀 Auth service starting on port %s", port)
	log.Printf("📖 Swagger UI: http://localhost:%s/swagger/", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
