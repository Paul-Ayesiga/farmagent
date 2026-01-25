package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/farmagent/fa-api-gateway/internal/config"
	"github.com/farmagent/fa-api-gateway/internal/middleware"
	"github.com/farmagent/fa-api-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
)

// @title FarmAgent API Gateway
// @version 1.0
// @description API Gateway for FarmAgent microservices - routes requests to auth, crop, payment, notification, and AI services.

// @contact.name FarmAgent Support
// @contact.email support@farmagent.ug

// @host localhost:8000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format: Bearer {token}

// proxyHealthTo creates a handler that proxies to a service's /health endpoint
func proxyHealthTo(serviceURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(serviceURL + "/health")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unavailable","error":"` + err.Error() + `"}`))
			return
		}
		defer resp.Body.Close()

		// Copy headers
		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}

func main() {
	// Load .env in development
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Load config
	cfg := config.Load()

	// Initialize proxy
	serviceProxy := proxy.NewServiceProxy()

	// Initialize router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rate limiting
	r.Use(httprate.LimitByIP(cfg.RateLimitRequests, time.Duration(cfg.RateLimitWindow)*time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","service":"fa-api-gateway"}`))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {

		// ===== Health Check Endpoints (Public) =====
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"healthy","service":"fa-api-gateway"}`))
		})
		r.Get("/health/auth", proxyHealthTo(cfg.AuthServiceURL))
		r.Get("/health/notification", proxyHealthTo(cfg.NotificationServiceURL))
		r.Get("/health/crop", proxyHealthTo(cfg.CropServiceURL))
		r.Get("/health/payment", proxyHealthTo(cfg.PaymentServiceURL))
		r.Get("/health/ai", proxyHealthTo(cfg.AIServiceURL))

		// ===== Auth Service (public routes) =====
		// Auth service routes are now mounted at /api/v1/auth, so we KEEP the prefix
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
			r.Post("/login", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
			r.Post("/refresh", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
			r.Post("/logout", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
			r.Post("/forgot-password", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
			r.Post("/reset-password", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
			r.Post("/verify-email", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
			r.Post("/resend-verification", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))

			// Protected auth routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuth(cfg))
				r.Get("/me", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
				r.Put("/me", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
				r.Post("/change-password", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
				r.Post("/send-verification", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
			})

			// Admin routes
			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuth(cfg))
				r.Use(middleware.RequireRole("admin"))
				r.Put("/users/{id}/role", serviceProxy.ProxyTo(cfg.AuthServiceURL, ""))
			})
		})

		// ===== Notification Service =====
		// Notification routes now mounted at /api/v1/notifications, KEEP prefix
		r.Route("/notifications", func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))
			r.Post("/email", serviceProxy.ProxyTo(cfg.NotificationServiceURL, ""))
			r.Post("/sms", serviceProxy.ProxyTo(cfg.NotificationServiceURL, ""))
			r.Post("/push", serviceProxy.ProxyTo(cfg.NotificationServiceURL, ""))
		})

		// ===== Crop Service (Protected) =====
		// Crop routes now mounted at /api/v1/..., KEEP prefix
		r.Route("/crops", func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))
			r.Get("/", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Post("/", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Get("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Put("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Delete("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			// Nested routes for health records and treatments
			r.Get("/{cropId}/health-records", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Get("/{cropId}/treatments", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
		})

		r.Route("/fields", func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))
			r.Get("/", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Post("/", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Get("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Put("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Delete("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Get("/{id}/crops", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
		})

		r.Route("/health-records", func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))
			r.Get("/", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Post("/", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Get("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Delete("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
		})

		r.Route("/treatments", func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))
			r.Post("/", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Get("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Put("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
			r.Delete("/{id}", serviceProxy.ProxyTo(cfg.CropServiceURL, ""))
		})

		// ===== AI Service (Protected) =====
		r.Route("/ai", func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))
			r.Post("/analyze", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			r.Post("/analyze/base64", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			r.Get("/analyze/status", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			r.Get("/analyze/diseases", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			r.Post("/recommend", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			r.Post("/chat", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			r.Post("/chat/stream", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			r.Get("/chat/suggestions", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			// Weather (public in AI service but protected via gateway)
			r.Get("/weather", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			r.Get("/weather/region/{region}", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			r.Get("/weather/regions", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
			r.Get("/spray-check", serviceProxy.ProxyTo(cfg.AIServiceURL, ""))
		})

		// ===== Payment Service (Protected) =====
		r.Route("/payments", func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))
			r.Post("/initiate", serviceProxy.ProxyTo(cfg.PaymentServiceURL, ""))
			r.Get("/{id}/status", serviceProxy.ProxyTo(cfg.PaymentServiceURL, ""))
			r.Get("/history", serviceProxy.ProxyTo(cfg.PaymentServiceURL, ""))
		})

		// Payment callback (NO AUTH - called by MTN)
		r.Post("/payments/callback", serviceProxy.ProxyTo(cfg.PaymentServiceURL, ""))
		r.Post("/payments/callback/{referenceId}", serviceProxy.ProxyTo(cfg.PaymentServiceURL, ""))

		// ===== Subscriptions (Protected) =====
		r.Route("/subscriptions", func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))
			r.Get("/", serviceProxy.ProxyTo(cfg.PaymentServiceURL, ""))
			r.Get("/plans", serviceProxy.ProxyTo(cfg.PaymentServiceURL, ""))
			r.Post("/subscribe", serviceProxy.ProxyTo(cfg.PaymentServiceURL, ""))
		})
	})

	// Start server
	port := cfg.AppPort
	if port == "" {
		port = "8000"
	}

	log.Printf("🚀 API Gateway starting on port %s", port)
	log.Printf("📍 Routes:")
	log.Printf("   /api/v1/health/*        → Service health checks")
	log.Printf("   /api/v1/auth/*          → %s", cfg.AuthServiceURL)
	log.Printf("   /api/v1/notifications/* → %s", cfg.NotificationServiceURL)
	log.Printf("   /api/v1/crops/*         → %s", cfg.CropServiceURL)
	log.Printf("   /api/v1/ai/*            → %s", cfg.AIServiceURL)
	log.Printf("   /api/v1/payments/*      → %s", cfg.PaymentServiceURL)
	log.Printf("   /api/v1/subscriptions/* → %s", cfg.PaymentServiceURL)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
