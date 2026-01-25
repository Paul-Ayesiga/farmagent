package main

import (
	"log"
	"net/http"
	"os"

	"github.com/farmagent/fa-crop-service/internal/config"
	"github.com/farmagent/fa-crop-service/internal/database"
	"github.com/farmagent/fa-crop-service/internal/handlers"
	"github.com/farmagent/fa-crop-service/internal/models"
	"github.com/farmagent/fa-crop-service/internal/repository"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

// @title FarmAgent Crop Service API
// @version 1.0
// @description Crop management service for FarmAgent - handles fields, crops, health records, and treatments.

// @host localhost:8002
// @BasePath /
// @schemes http https

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

	// Auto-migrate
	log.Println("Running migrations...")
	db.AutoMigrate(&models.Field{}, &models.Crop{}, &models.HealthRecord{}, &models.Treatment{})

	// Initialize repositories
	fieldRepo := repository.NewFieldRepository(db)
	cropRepo := repository.NewCropRepository(db)
	healthRepo := repository.NewHealthRecordRepository(db)
	treatmentRepo := repository.NewTreatmentRepository(db)

	// Initialize handlers
	fieldHandler := handlers.NewFieldHandler(fieldRepo)
	cropHandler := handlers.NewCropHandler(cropRepo, fieldRepo)
	healthHandler := handlers.NewHealthRecordHandler(healthRepo, cropRepo, fieldRepo)
	treatmentHandler := handlers.NewTreatmentHandler(treatmentRepo, cropRepo, fieldRepo)

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-User-ID", "X-User-Role"},
		AllowCredentials: true,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","service":"fa-crop-service"}`))
	})

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Field routes
		r.Route("/fields", func(r chi.Router) {
			r.Get("/", fieldHandler.ListFields)
			r.Post("/", fieldHandler.CreateField)
			r.Get("/{id}", fieldHandler.GetField)
			r.Put("/{id}", fieldHandler.UpdateField)
			r.Delete("/{id}", fieldHandler.DeleteField)
			r.Get("/{fieldId}/crops", cropHandler.ListCropsByField)
		})

		// Crop routes
		r.Route("/crops", func(r chi.Router) {
			r.Get("/", cropHandler.ListCrops)
			r.Post("/", cropHandler.CreateCrop)
			r.Get("/{id}", cropHandler.GetCrop)
			r.Put("/{id}", cropHandler.UpdateCrop)
			r.Delete("/{id}", cropHandler.DeleteCrop)
			r.Get("/{cropId}/health-records", healthHandler.ListHealthRecords)
			r.Get("/{cropId}/treatments", treatmentHandler.ListTreatments)
		})

		// Health record routes
		r.Route("/health-records", func(r chi.Router) {
			r.Post("/", healthHandler.CreateHealthRecord)
			r.Get("/{id}", healthHandler.GetHealthRecord)
			r.Delete("/{id}", healthHandler.DeleteHealthRecord)
		})

		// Treatment routes
		r.Route("/treatments", func(r chi.Router) {
			r.Post("/", treatmentHandler.CreateTreatment)
			r.Get("/{id}", treatmentHandler.GetTreatment)
			r.Put("/{id}", treatmentHandler.UpdateTreatment)
			r.Delete("/{id}", treatmentHandler.DeleteTreatment)
		})
	})

	// Start server
	port := cfg.AppPort
	if port == "" {
		port = "8002"
	}

	log.Printf("🌱 Crop service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
