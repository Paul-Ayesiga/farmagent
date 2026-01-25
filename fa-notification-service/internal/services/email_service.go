package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"strconv"

	"github.com/farmagent/fa-notification-service/internal/config"
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	Send(to, subject, body string) error
	SendTemplate(to, subject, template string, data map[string]string) error
}

type emailService struct {
	cfg    *config.Config
	dialer *gomail.Dialer
}

func NewEmailService(cfg *config.Config) EmailService {
	port, _ := strconv.Atoi(cfg.SMTPPort)
	dialer := gomail.NewDialer(cfg.SMTPHost, port, cfg.SMTPUser, cfg.SMTPPassword)

	// For Mailpit (no TLS in dev)
	if cfg.AppEnv == "development" {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &emailService{cfg: cfg, dialer: dialer}
}

func (s *emailService) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.cfg.SMTPFromName, s.cfg.SMTPFrom))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("✉️ Email sent to %s: %s", to, subject)
	return nil
}

func (s *emailService) SendTemplate(to, subject, template string, data map[string]string) error {
	body := s.renderTemplate(template, data)
	return s.Send(to, subject, body)
}

func (s *emailService) renderTemplate(template string, data map[string]string) string {
	// Simple template rendering - replace with proper templating in production
	templates := map[string]string{
		"welcome": `
			<html>
			<body style="font-family: Arial, sans-serif; padding: 20px;">
				<h1>Welcome to FarmAgent, {{name}}! 🌱</h1>
				<p>We're excited to have you on board.</p>
				<p>Start by registering your fields and crops to get personalized recommendations.</p>
				<p>Best regards,<br>The FarmAgent Team</p>
			</body>
			</html>
		`,
		"password_reset": `
			<html>
			<body style="font-family: Arial, sans-serif; padding: 20px;">
				<h2>Password Reset Request</h2>
				<p>Hello {{name}},</p>
				<p>You requested to reset your password. Use this code:</p>
				<h1 style="background: #f0f0f0; padding: 20px; text-align: center;">{{code}}</h1>
				<p>This code expires in 1 hour.</p>
				<p>If you didn't request this, please ignore this email.</p>
			</body>
			</html>
		`,
		"disease_alert": `
			<html>
			<body style="font-family: Arial, sans-serif; padding: 20px;">
				<h2>⚠️ Disease Alert for Your Crop</h2>
				<p>Hello {{name}},</p>
				<p>We detected <strong>{{disease}}</strong> in your {{crop}}.</p>
				<p><strong>Recommended Treatment:</strong> {{treatment}}</p>
				<p>Take action soon to minimize damage.</p>
			</body>
			</html>
		`,
		"verification": `
			<html>
			<body style="font-family: Arial, sans-serif; padding: 20px;">
				<h2>Verify Your Email</h2>
				<p>Hello {{name}},</p>
				<p>Your verification code is:</p>
				<h1 style="background: #4CAF50; color: white; padding: 20px; text-align: center;">{{code}}</h1>
				<p>This code expires in 24 hours.</p>
			</body>
			</html>
		`,
		"email_verification": `
			<html>
			<body style="font-family: Arial, sans-serif; padding: 20px;">
				<h2>Verify Your FarmAgent Email</h2>
				<p>Hello {{name}},</p>
				<p>Your email verification code is:</p>
				<h1 style="background: #4CAF50; color: white; padding: 20px; text-align: center; font-size: 32px; letter-spacing: 8px;">{{code}}</h1>
				<p>This code expires in 24 hours.</p>
				<p>If you didn't create an account with FarmAgent, please ignore this email.</p>
				<br>
				<p>Best regards,<br>The FarmAgent Team</p>
			</body>
			</html>
		`,
	}

	body, ok := templates[template]
	if !ok {
		return "Template not found"
	}

	// Simple placeholder replacement
	for key, value := range data {
		body = replaceAll(body, "{{"+key+"}}", value)
	}

	return body
}

func replaceAll(s, old, new string) string {
	for {
		i := indexOf(s, old)
		if i == -1 {
			return s
		}
		s = s[:i] + new + s[i+len(old):]
	}
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
