package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/farmagent/fa-notification-service/internal/queue"
)

type NotificationHandler struct {
	rabbitmq *queue.RabbitMQ
}

func NewNotificationHandler(rabbitmq *queue.RabbitMQ) *NotificationHandler {
	return &NotificationHandler{rabbitmq: rabbitmq}
}

type SendEmailRequest struct {
	To       string            `json:"to"`
	Subject  string            `json:"subject"`
	Template string            `json:"template,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

type SendSMSRequest struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type SendPushRequest struct {
	DeviceToken string            `json:"device_token"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Data        map[string]string `json:"data,omitempty"`
}

// SendEmail queues an email for background sending
func (h *NotificationHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.To == "" || req.Subject == "" {
		respondError(w, http.StatusBadRequest, "to and subject are required")
		return
	}

	msg := queue.EmailMessage{
		To:       req.To,
		Subject:  req.Subject,
		Template: req.Template,
		Data:     req.Data,
	}

	if err := h.rabbitmq.PublishEmail(context.Background(), msg); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to queue email")
		return
	}

	respondJSON(w, http.StatusAccepted, map[string]string{
		"message": "email queued successfully",
	})
}

// SendSMS queues an SMS for background sending
func (h *NotificationHandler) SendSMS(w http.ResponseWriter, r *http.Request) {
	var req SendSMSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Phone == "" || req.Message == "" {
		respondError(w, http.StatusBadRequest, "phone and message are required")
		return
	}

	msg := queue.SMSMessage{
		Phone:   req.Phone,
		Message: req.Message,
	}

	if err := h.rabbitmq.PublishSMS(context.Background(), msg); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to queue SMS")
		return
	}

	respondJSON(w, http.StatusAccepted, map[string]string{
		"message": "sms queued successfully",
	})
}

// SendPush queues a push notification
func (h *NotificationHandler) SendPush(w http.ResponseWriter, r *http.Request) {
	var req SendPushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DeviceToken == "" || req.Title == "" {
		respondError(w, http.StatusBadRequest, "device_token and title are required")
		return
	}

	msg := queue.PushMessage{
		DeviceToken: req.DeviceToken,
		Title:       req.Title,
		Body:        req.Body,
		Data:        req.Data,
	}

	if err := h.rabbitmq.PublishPush(context.Background(), msg); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to queue push notification")
		return
	}

	respondJSON(w, http.StatusAccepted, map[string]string{
		"message": "push notification queued successfully",
	})
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
