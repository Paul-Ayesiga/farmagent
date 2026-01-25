package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/farmagent/fa-payment-service/internal/models"
	"github.com/farmagent/fa-payment-service/internal/repository"
)

// PaymentService handles payment business logic
type PaymentService struct {
	mtn     *MTNMoMoClient
	txRepo  *repository.TransactionRepository
	subRepo *repository.SubscriptionRepository
}

// NewPaymentService creates a new payment service
func NewPaymentService(mtn *MTNMoMoClient, txRepo *repository.TransactionRepository, subRepo *repository.SubscriptionRepository) *PaymentService {
	return &PaymentService{
		mtn:     mtn,
		txRepo:  txRepo,
		subRepo: subRepo,
	}
}

// InitiatePaymentRequest represents a payment initiation request
type InitiatePaymentRequest struct {
	UserID uuid.UUID
	Amount float64
	Phone  string
	Reason string
	PlanID string // optional, for subscription payments
}

// InitiatePaymentResponse represents the result of initiating a payment
type InitiatePaymentResponse struct {
	TransactionID  uuid.UUID `json:"transaction_id"`
	MTNReferenceID string    `json:"mtn_reference_id"`
	Status         string    `json:"status"`
	Message        string    `json:"message"`
}

// InitiatePayment starts a new payment collection
func (s *PaymentService) InitiatePayment(req InitiatePaymentRequest) (*InitiatePaymentResponse, error) {
	// Generate IDs
	txID := uuid.New()
	externalID := fmt.Sprintf("FA-%s", txID.String()[:8])

	// Create transaction record first (pending)
	tx := &models.Transaction{
		ID:           txID,
		UserID:       req.UserID,
		ExternalID:   externalID,
		Type:         models.TypeCollection,
		Amount:       req.Amount,
		Currency:     "UGX",
		Phone:        req.Phone,
		Status:       models.StatusPending,
		PayerMessage: req.Reason,
		PayeeNote:    "FarmAgent Payment",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.txRepo.Create(tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Create pending subscription if this is for a plan
	if req.PlanID != "" {
		plan, ok := models.GetPlan(req.PlanID)
		if !ok {
			return nil, fmt.Errorf("invalid plan ID: %s", req.PlanID)
		}

		sub := &models.Subscription{
			ID:        uuid.New(),
			UserID:    req.UserID,
			PlanID:    req.PlanID,
			Status:    models.SubStatusPending,
			AutoRenew: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			// Start/End dates set after payment confirmation
			StartDate: time.Time{},
			EndDate:   time.Time{},
		}

		if err := s.subRepo.Create(sub); err != nil {
			return nil, fmt.Errorf("failed to create subscription: %w", err)
		}

		// Use plan amount if not specified
		if req.Amount == 0 {
			req.Amount = plan.Amount
		}
	}

	// Call MTN MoMo API
	referenceID, err := s.mtn.RequestToPay(
		req.Amount,
		"UGX",
		req.Phone,
		externalID,
		req.Reason,
		"FarmAgent Payment",
	)
	if err != nil {
		// Update transaction as failed
		s.txRepo.UpdateStatus(txID, models.StatusFailed, err.Error(), nil)
		return nil, fmt.Errorf("MTN request failed: %w", err)
	}

	// Update transaction with MTN reference
	tx.MTNReferenceID = referenceID
	s.txRepo.UpdateStatus(txID, models.StatusPending, "", nil)

	return &InitiatePaymentResponse{
		TransactionID:  txID,
		MTNReferenceID: referenceID,
		Status:         models.StatusPending,
		Message:        "Payment request sent. Please approve on your phone.",
	}, nil
}

// CheckTransactionStatus checks status from MTN and updates local record
func (s *PaymentService) CheckTransactionStatus(txID uuid.UUID) (*models.Transaction, error) {
	// Get local transaction
	tx, err := s.txRepo.GetByID(txID)
	if err != nil {
		return nil, fmt.Errorf("transaction not found: %w", err)
	}

	// If already final, return as-is
	if tx.Status == models.StatusSuccessful || tx.Status == models.StatusFailed {
		return tx, nil
	}

	// Check status from MTN
	result, err := s.mtn.GetTransactionStatus(tx.MTNReferenceID)
	if err != nil {
		return nil, fmt.Errorf("failed to check MTN status: %w", err)
	}

	// Map MTN status to our status
	var newStatus, failureReason string
	switch result.Status {
	case "SUCCESSFUL":
		newStatus = models.StatusSuccessful
	case "FAILED":
		newStatus = models.StatusFailed
		if result.Reason != nil {
			failureReason = result.Reason.Message
		}
	case "PENDING":
		newStatus = models.StatusPending
	default:
		newStatus = models.StatusPending
	}

	// Update if changed
	if newStatus != tx.Status {
		callbackPayload := models.JSONMap{
			"status":                 result.Status,
			"amount":                 result.Amount,
			"financialTransactionId": result.FinancialTransactionID,
		}
		s.txRepo.UpdateStatus(txID, newStatus, failureReason, callbackPayload)
		tx.Status = newStatus
		tx.FailureReason = failureReason

		// If successful and it's a subscription payment, activate subscription
		if newStatus == models.StatusSuccessful {
			s.activateSubscriptionForUser(tx.UserID, txID)
		}
	}

	return tx, nil
}

// ProcessCallback processes MTN callback
func (s *PaymentService) ProcessCallback(referenceID string, payload map[string]interface{}) error {
	// Find transaction by MTN reference
	tx, err := s.txRepo.GetByMTNReferenceID(referenceID)
	if err != nil {
		return fmt.Errorf("transaction not found for reference %s: %w", referenceID, err)
	}

	// Extract status from payload
	status, ok := payload["status"].(string)
	if !ok {
		return fmt.Errorf("invalid callback payload: missing status")
	}

	var newStatus, failureReason string
	switch status {
	case "SUCCESSFUL":
		newStatus = models.StatusSuccessful
	case "FAILED":
		newStatus = models.StatusFailed
		if reason, ok := payload["reason"].(map[string]interface{}); ok {
			if msg, ok := reason["message"].(string); ok {
				failureReason = msg
			}
		}
	default:
		return nil // Ignore other statuses
	}

	// Update transaction
	s.txRepo.UpdateStatus(tx.ID, newStatus, failureReason, models.JSONMap(payload))

	// If successful, activate subscription
	if newStatus == models.StatusSuccessful {
		s.activateSubscriptionForUser(tx.UserID, tx.ID)
	}

	return nil
}

// activateSubscriptionForUser activates pending subscription after payment
func (s *PaymentService) activateSubscriptionForUser(userID, transactionID uuid.UUID) error {
	// Find pending subscription for this user
	sub, err := s.subRepo.GetByUserID(userID)
	if err != nil {
		// No subscription to activate, maybe a one-time payment
		return nil
	}

	if sub.Status != models.SubStatusPending && sub.Status != models.SubStatusActive {
		return nil
	}

	plan, ok := models.GetPlan(sub.PlanID)
	if !ok {
		return fmt.Errorf("invalid plan: %s", sub.PlanID)
	}

	startDate := time.Now()
	endDate := startDate.AddDate(0, 0, plan.Duration)

	return s.subRepo.ActivateSubscription(sub.ID, transactionID, startDate, endDate)
}

// GetUserTransactions returns transaction history for a user
func (s *PaymentService) GetUserTransactions(userID uuid.UUID, limit, offset int) ([]*models.Transaction, error) {
	return s.txRepo.GetByUserID(userID, limit, offset)
}

// GetUserSubscription returns the active subscription for a user
func (s *PaymentService) GetUserSubscription(userID uuid.UUID) (*models.Subscription, error) {
	return s.subRepo.GetByUserID(userID)
}
