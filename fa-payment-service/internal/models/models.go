package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Transaction status constants
const (
	StatusPending    = "pending"
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
	StatusCancelled  = "cancelled"
)

// Transaction type constants
const (
	TypeCollection   = "collection"
	TypeDisbursement = "disbursement"
)

// JSONMap for storing callback payload
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Transaction represents a payment transaction
type Transaction struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	ExternalID      string    `json:"external_id"`      // Our reference
	MTNReferenceID  string    `json:"mtn_reference_id"` // MTN's X-Reference-Id
	Type            string    `json:"type"`             // collection, disbursement
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Phone           string    `json:"phone"`
	Status          string    `json:"status"`
	FailureReason   string    `json:"failure_reason,omitempty"`
	PayerMessage    string    `json:"payer_message,omitempty"`
	PayeeNote       string    `json:"payee_note,omitempty"`
	CallbackPayload JSONMap   `json:"callback_payload,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Subscription status constants
const (
	SubStatusActive    = "active"
	SubStatusCancelled = "cancelled"
	SubStatusExpired   = "expired"
	SubStatusPending   = "pending"
)

// Subscription represents a user's subscription
type Subscription struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	PlanID        string     `json:"plan_id"`
	TransactionID *uuid.UUID `json:"transaction_id,omitempty"`
	Status        string     `json:"status"`
	StartDate     time.Time  `json:"start_date"`
	EndDate       time.Time  `json:"end_date"`
	AutoRenew     bool       `json:"auto_renew"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Plan represents a subscription plan
type Plan struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Amount      float64  `json:"amount"`
	Currency    string   `json:"currency"`
	Duration    int      `json:"duration"` // in days
	Features    []string `json:"features"`
}

// Available plans
var Plans = map[string]Plan{
	"free": {
		ID:          "free",
		Name:        "Free Plan",
		Description: "Basic access to FarmAgent",
		Amount:      0,
		Currency:    "UGX",
		Duration:    0, // unlimited
		Features: []string{
			"2 disease scans per month",
			"Basic recommendations",
		},
	},
	"basic": {
		ID:          "basic",
		Name:        "Basic Plan",
		Description: "Essential features for farmers",
		Amount:      5000, // 5,000 UGX
		Currency:    "UGX",
		Duration:    30,
		Features: []string{
			"10 disease scans per month",
			"Treatment recommendations",
			"Weather alerts",
		},
	},
	"premium": {
		ID:          "premium",
		Name:        "Premium Plan",
		Description: "Full access to all FarmAgent features",
		Amount:      20000, // 20,000 UGX
		Currency:    "UGX",
		Duration:    30,
		Features: []string{
			"Unlimited disease scans",
			"AI chat assistant",
			"Weather alerts",
			"Market price updates",
			"Expert consultations",
			"Priority support",
		},
	},
}

// GetPlan returns a plan by ID
func GetPlan(planID string) (Plan, bool) {
	plan, ok := Plans[planID]
	return plan, ok
}

// GetAllPlans returns all available plans
func GetAllPlans() []Plan {
	plans := make([]Plan, 0, len(Plans))
	for _, p := range Plans {
		plans = append(plans, p)
	}
	return plans
}
