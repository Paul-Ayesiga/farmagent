package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/farmagent/fa-payment-service/internal/config"
)

// MTNMoMoClient handles communication with MTN MoMo API
type MTNMoMoClient struct {
	cfg         *config.Config
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
	tokenMutex  sync.RWMutex
}

// NewMTNMoMoClient creates a new MTN MoMo API client
func NewMTNMoMoClient(cfg *config.Config) *MTNMoMoClient {
	return &MTNMoMoClient{
		cfg: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RequestToPayRequest is the request body for initiating a collection
type RequestToPayRequest struct {
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	ExternalID   string `json:"externalId"`
	Payer        Payer  `json:"payer"`
	PayerMessage string `json:"payerMessage"`
	PayeeNote    string `json:"payeeNote"`
}

// Payer represents the payer information
type Payer struct {
	PartyIDType string `json:"partyIdType"`
	PartyID     string `json:"partyId"`
}

// RequestToPayResponse is the response from checking transaction status
type RequestToPayResponse struct {
	Amount                 string `json:"amount"`
	Currency               string `json:"currency"`
	FinancialTransactionID string `json:"financialTransactionId"`
	ExternalID             string `json:"externalId"`
	Payer                  Payer  `json:"payer"`
	PayerMessage           string `json:"payerMessage"`
	PayeeNote              string `json:"payeeNote"`
	Status                 string `json:"status"` // PENDING, SUCCESSFUL, FAILED
	Reason                 *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"reason,omitempty"`
}

// TokenResponse is the response from token endpoint
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// GetAccessToken retrieves or refreshes the access token
func (c *MTNMoMoClient) GetAccessToken() (string, error) {
	c.tokenMutex.RLock()
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		defer c.tokenMutex.RUnlock()
		return c.accessToken, nil
	}
	c.tokenMutex.RUnlock()

	// Need to refresh token
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	// Double-check after acquiring write lock
	if c.accessToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	// Create Basic Auth header
	credentials := fmt.Sprintf("%s:%s", c.cfg.MTNAPIUser, c.cfg.MTNAPIKey)
	basicAuth := base64.StdEncoding.EncodeToString([]byte(credentials))

	url := fmt.Sprintf("%s/collection/token/", c.cfg.MTNBaseURL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Authorization", "Basic "+basicAuth)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.cfg.MTNSubscriptionKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	// Set expiry with some buffer (subtract 60 seconds)
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn-60) * time.Second)

	return c.accessToken, nil
}

// RequestToPay initiates a collection request
func (c *MTNMoMoClient) RequestToPay(amount float64, currency, phone, externalID, message, note string) (string, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}

	referenceID := uuid.New().String()

	reqBody := RequestToPayRequest{
		Amount:       fmt.Sprintf("%.0f", amount),
		Currency:     currency,
		ExternalID:   externalID,
		PayerMessage: message,
		PayeeNote:    note,
		Payer: Payer{
			PartyIDType: "MSISDN",
			PartyID:     phone,
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/collection/v1_0/requesttopay", c.cfg.MTNBaseURL)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Reference-Id", referenceID)
	req.Header.Set("X-Target-Environment", c.cfg.MTNEnvironment)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.cfg.MTNSubscriptionKey)
	req.Header.Set("Content-Type", "application/json")

	if c.cfg.MTNCallbackURL != "" {
		req.Header.Set("X-Callback-Url", c.cfg.MTNCallbackURL)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request to pay failed: %w", err)
	}
	defer resp.Body.Close()

	// 202 Accepted means the request was accepted for processing
	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request to pay failed with status %d: %s", resp.StatusCode, string(body))
	}

	return referenceID, nil
}

// GetTransactionStatus checks the status of a transaction
func (c *MTNMoMoClient) GetTransactionStatus(referenceID string) (*RequestToPayResponse, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/collection/v1_0/requesttopay/%s", c.cfg.MTNBaseURL, referenceID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Target-Environment", c.cfg.MTNEnvironment)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.cfg.MTNSubscriptionKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get transaction status failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get transaction status failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result RequestToPayResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetAccountBalance returns the account balance
func (c *MTNMoMoClient) GetAccountBalance() (map[string]interface{}, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/collection/v1_0/account/balance", c.cfg.MTNBaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Target-Environment", c.cfg.MTNEnvironment)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.cfg.MTNSubscriptionKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get balance failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get balance failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// ValidateAccountHolder validates if an account exists
func (c *MTNMoMoClient) ValidateAccountHolder(phone string) (bool, error) {
	token, err := c.GetAccessToken()
	if err != nil {
		return false, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/collection/v1_0/accountholder/msisdn/%s/active", c.cfg.MTNBaseURL, phone)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Target-Environment", c.cfg.MTNEnvironment)
	req.Header.Set("Ocp-Apim-Subscription-Key", c.cfg.MTNSubscriptionKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("validate account failed: %w", err)
	}
	defer resp.Body.Close()

	// 200 means account exists and is active
	return resp.StatusCode == http.StatusOK, nil
}
