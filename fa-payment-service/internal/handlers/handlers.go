package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/farmagent/fa-payment-service/internal/models"
	"github.com/farmagent/fa-payment-service/internal/services"
)

// PaymentHandler handles payment HTTP requests
type PaymentHandler struct {
	paymentSvc *services.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(paymentSvc *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentSvc: paymentSvc}
}

// InitiatePaymentRequest is the request body for initiating payment
type InitiatePaymentRequest struct {
	Amount float64 `json:"amount"`
	Phone  string  `json:"phone"`
	Reason string  `json:"reason"`
	PlanID string  `json:"plan_id,omitempty"`
}

// InitiatePayment handles POST /payments/initiate
func (h *PaymentHandler) InitiatePayment(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req InitiatePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate phone number (Uganda format)
	if len(req.Phone) < 10 {
		writeError(w, http.StatusBadRequest, "invalid phone number")
		return
	}

	// If plan_id provided, use plan amount
	if req.PlanID != "" {
		plan, ok := models.GetPlan(req.PlanID)
		if !ok {
			writeError(w, http.StatusBadRequest, "invalid plan ID")
			return
		}
		if req.Amount == 0 {
			req.Amount = plan.Amount
		}
		if req.Reason == "" {
			req.Reason = plan.Name + " subscription"
		}
	}

	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}

	result, err := h.paymentSvc.InitiatePayment(services.InitiatePaymentRequest{
		UserID: userID,
		Amount: req.Amount,
		Phone:  req.Phone,
		Reason: req.Reason,
		PlanID: req.PlanID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, result)
}

// GetTransactionStatus handles GET /payments/{id}/status
func (h *PaymentHandler) GetTransactionStatus(w http.ResponseWriter, r *http.Request) {
	txIDStr := chi.URLParam(r, "id")
	txID, err := uuid.Parse(txIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid transaction ID")
		return
	}

	tx, err := h.paymentSvc.CheckTransactionStatus(txID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tx)
}

// Callback handles POST /payments/callback (from MTN)
func (h *PaymentHandler) Callback(w http.ResponseWriter, r *http.Request) {
	// Get reference ID from header or path
	referenceID := r.Header.Get("X-Reference-Id")
	if referenceID == "" {
		referenceID = chi.URLParam(r, "referenceId")
	}

	if referenceID == "" {
		writeError(w, http.StatusBadRequest, "missing reference ID")
		return
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "invalid callback payload")
		return
	}

	if err := h.paymentSvc.ProcessCallback(referenceID, payload); err != nil {
		// Log error but return 200 to MTN (don't want retries for our errors)
		// In production, log this properly
		writeJSON(w, http.StatusOK, map[string]string{"status": "processed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// GetHistory handles GET /payments/history
func (h *PaymentHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Parse pagination
	limit := 10
	offset := 0
	// Could parse from query params

	transactions, err := h.paymentSvc.GetUserTransactions(userID, limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"transactions": transactions,
		"count":        len(transactions),
	})
}

// SubscriptionHandler handles subscription HTTP requests
type SubscriptionHandler struct {
	paymentSvc *services.PaymentService
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(paymentSvc *services.PaymentService) *SubscriptionHandler {
	return &SubscriptionHandler{paymentSvc: paymentSvc}
}

// GetPlans handles GET /subscriptions/plans
func (h *SubscriptionHandler) GetPlans(w http.ResponseWriter, r *http.Request) {
	plans := models.GetAllPlans()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"plans": plans,
	})
}

// GetCurrentSubscription handles GET /subscriptions
func (h *SubscriptionHandler) GetCurrentSubscription(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	sub, err := h.paymentSvc.GetUserSubscription(userID)
	if err != nil {
		// No subscription found - return free plan info
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"subscription": nil,
			"current_plan": models.Plans["free"],
			"message":      "No active subscription. Using free plan.",
		})
		return
	}

	plan, _ := models.GetPlan(sub.PlanID)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"subscription": sub,
		"current_plan": plan,
	})
}

// Subscribe handles POST /subscriptions/subscribe
func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		writeError(w, http.StatusUnauthorized, "missing user ID")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req struct {
		PlanID string `json:"plan_id"`
		Phone  string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	plan, ok := models.GetPlan(req.PlanID)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid plan ID")
		return
	}

	if plan.Amount == 0 {
		writeError(w, http.StatusBadRequest, "cannot subscribe to free plan via payment")
		return
	}

	// Initiate payment for subscription
	result, err := h.paymentSvc.InitiatePayment(services.InitiatePaymentRequest{
		UserID: userID,
		Amount: plan.Amount,
		Phone:  req.Phone,
		Reason: plan.Name + " - FarmAgent",
		PlanID: req.PlanID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]interface{}{
		"transaction": result,
		"plan":        plan,
		"message":     "Payment request sent. Approve on your phone to activate subscription.",
	})
}

// Helper functions

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
		},
	})
}
