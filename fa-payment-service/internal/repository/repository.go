package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/farmagent/fa-payment-service/internal/models"
)

// TransactionRepository handles transaction database operations
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create inserts a new transaction
func (r *TransactionRepository) Create(tx *models.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, user_id, external_id, mtn_reference_id, type, amount, currency, 
			phone, status, payer_message, payee_note, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.Exec(query,
		tx.ID, tx.UserID, tx.ExternalID, tx.MTNReferenceID, tx.Type,
		tx.Amount, tx.Currency, tx.Phone, tx.Status,
		tx.PayerMessage, tx.PayeeNote, tx.CreatedAt, tx.UpdatedAt,
	)
	return err
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(id uuid.UUID) (*models.Transaction, error) {
	query := `
		SELECT id, user_id, external_id, mtn_reference_id, type, amount, currency,
			   phone, status, failure_reason, payer_message, payee_note, 
			   callback_payload, created_at, updated_at
		FROM transactions WHERE id = $1
	`

	tx := &models.Transaction{}
	err := r.db.QueryRow(query, id).Scan(
		&tx.ID, &tx.UserID, &tx.ExternalID, &tx.MTNReferenceID, &tx.Type,
		&tx.Amount, &tx.Currency, &tx.Phone, &tx.Status, &tx.FailureReason,
		&tx.PayerMessage, &tx.PayeeNote, &tx.CallbackPayload, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetByMTNReferenceID retrieves a transaction by MTN reference ID
func (r *TransactionRepository) GetByMTNReferenceID(refID string) (*models.Transaction, error) {
	query := `
		SELECT id, user_id, external_id, mtn_reference_id, type, amount, currency,
			   phone, status, failure_reason, payer_message, payee_note,
			   callback_payload, created_at, updated_at
		FROM transactions WHERE mtn_reference_id = $1
	`

	tx := &models.Transaction{}
	err := r.db.QueryRow(query, refID).Scan(
		&tx.ID, &tx.UserID, &tx.ExternalID, &tx.MTNReferenceID, &tx.Type,
		&tx.Amount, &tx.Currency, &tx.Phone, &tx.Status, &tx.FailureReason,
		&tx.PayerMessage, &tx.PayeeNote, &tx.CallbackPayload, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetByUserID retrieves all transactions for a user
func (r *TransactionRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*models.Transaction, error) {
	query := `
		SELECT id, user_id, external_id, mtn_reference_id, type, amount, currency,
			   phone, status, failure_reason, payer_message, payee_note,
			   callback_payload, created_at, updated_at
		FROM transactions WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		tx := &models.Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.ExternalID, &tx.MTNReferenceID, &tx.Type,
			&tx.Amount, &tx.Currency, &tx.Phone, &tx.Status, &tx.FailureReason,
			&tx.PayerMessage, &tx.PayeeNote, &tx.CallbackPayload, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

// UpdateStatus updates transaction status and failure reason
func (r *TransactionRepository) UpdateStatus(id uuid.UUID, status, failureReason string, callbackPayload models.JSONMap) error {
	query := `
		UPDATE transactions 
		SET status = $2, failure_reason = $3, callback_payload = $4, updated_at = $5
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id, status, failureReason, callbackPayload, time.Now())
	return err
}

// SubscriptionRepository handles subscription database operations
type SubscriptionRepository struct {
	db *sql.DB
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create inserts a new subscription
func (r *SubscriptionRepository) Create(sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (
			id, user_id, plan_id, transaction_id, status, start_date, end_date,
			auto_renew, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(query,
		sub.ID, sub.UserID, sub.PlanID, sub.TransactionID, sub.Status,
		sub.StartDate, sub.EndDate, sub.AutoRenew, sub.CreatedAt, sub.UpdatedAt,
	)
	return err
}

// GetByUserID retrieves active subscription for a user
func (r *SubscriptionRepository) GetByUserID(userID uuid.UUID) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, plan_id, transaction_id, status, start_date, end_date,
			   auto_renew, created_at, updated_at
		FROM subscriptions 
		WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1
	`

	sub := &models.Subscription{}
	err := r.db.QueryRow(query, userID).Scan(
		&sub.ID, &sub.UserID, &sub.PlanID, &sub.TransactionID, &sub.Status,
		&sub.StartDate, &sub.EndDate, &sub.AutoRenew, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// UpdateStatus updates subscription status
func (r *SubscriptionRepository) UpdateStatus(id uuid.UUID, status string) error {
	query := `UPDATE subscriptions SET status = $2, updated_at = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, status, time.Now())
	return err
}

// ActivateSubscription activates a subscription after successful payment
func (r *SubscriptionRepository) ActivateSubscription(id uuid.UUID, transactionID uuid.UUID, startDate, endDate time.Time) error {
	query := `
		UPDATE subscriptions 
		SET status = 'active', transaction_id = $2, start_date = $3, end_date = $4, updated_at = $5
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, transactionID, startDate, endDate, time.Now())
	return err
}
