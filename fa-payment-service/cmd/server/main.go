package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	"github.com/farmagent/fa-payment-service/internal/config"
	"github.com/farmagent/fa-payment-service/internal/handlers"
	"github.com/farmagent/fa-payment-service/internal/repository"
	"github.com/farmagent/fa-payment-service/internal/services"
)

func main() {
	cfg := config.Load()

	log.Printf("💳 Starting %s on port %s", cfg.AppName, cfg.AppPort)
	log.Printf("   Environment: %s", cfg.AppEnv)
	log.Printf("   MTN Environment: %s", cfg.MTNEnvironment)

	// Database connection
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	mtnClient := services.NewMTNMoMoClient(cfg)
	txRepo := repository.NewTransactionRepository(db)
	subRepo := repository.NewSubscriptionRepository(db)
	paymentSvc := services.NewPaymentService(mtnClient, txRepo, subRepo)

	// Initialize handlers
	paymentHandler := handlers.NewPaymentHandler(paymentSvc)
	subscriptionHandler := handlers.NewSubscriptionHandler(paymentSvc)

	// Router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","service":"fa-payment-service"}`))
	})

	// MTN Callback (no auth required - called by MTN)
	r.Post("/api/v1/payments/callback", paymentHandler.Callback)
	r.Post("/api/v1/payments/callback/{referenceId}", paymentHandler.Callback)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Payment routes
		r.Route("/payments", func(r chi.Router) {
			r.Post("/initiate", paymentHandler.InitiatePayment)
			r.Get("/{id}/status", paymentHandler.GetTransactionStatus)
			r.Get("/history", paymentHandler.GetHistory)
		})

		// Subscription routes
		r.Route("/subscriptions", func(r chi.Router) {
			r.Get("/", subscriptionHandler.GetCurrentSubscription)
			r.Get("/plans", subscriptionHandler.GetPlans)
			r.Post("/subscribe", subscriptionHandler.Subscribe)
		})
	})

	// Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("🚀 Server listening on %s", addr)
	log.Printf("📡 MTN Callback URL: %s", cfg.MTNCallbackURL)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func initDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("✅ Connected to database")
	return db, nil
}

func runMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS transactions (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			external_id VARCHAR(100) UNIQUE NOT NULL,
			mtn_reference_id VARCHAR(100),
			type VARCHAR(50) NOT NULL,
			amount DECIMAL(15,2) NOT NULL,
			currency VARCHAR(10) NOT NULL DEFAULT 'UGX',
			phone VARCHAR(20) NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			failure_reason TEXT,
			payer_message TEXT,
			payee_note TEXT,
			callback_payload JSONB,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_mtn_reference ON transactions(mtn_reference_id)`,
		`CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status)`,

		`CREATE TABLE IF NOT EXISTS subscriptions (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL,
			plan_id VARCHAR(50) NOT NULL,
			transaction_id UUID REFERENCES transactions(id),
			status VARCHAR(50) NOT NULL DEFAULT 'pending',
			start_date TIMESTAMP,
			end_date TIMESTAMP,
			auto_renew BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status)`,
	}

	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	log.Println("✅ Database migrations complete")
	return nil
}
