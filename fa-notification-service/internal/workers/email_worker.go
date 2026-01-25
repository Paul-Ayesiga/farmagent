package workers

import (
	"log"

	"github.com/farmagent/fa-notification-service/internal/queue"
	"github.com/farmagent/fa-notification-service/internal/services"
)

type EmailWorker struct {
	emailService services.EmailService
	rabbitmq     *queue.RabbitMQ
}

func NewEmailWorker(emailService services.EmailService, rabbitmq *queue.RabbitMQ) *EmailWorker {
	return &EmailWorker{
		emailService: emailService,
		rabbitmq:     rabbitmq,
	}
}

func (w *EmailWorker) Start() error {
	log.Println("📧 Email worker started, listening for messages...")

	return w.rabbitmq.ConsumeEmails(func(msg queue.EmailMessage) error {
		log.Printf("Processing email to: %s, template: %s", msg.To, msg.Template)

		if msg.Template != "" {
			return w.emailService.SendTemplate(msg.To, msg.Subject, msg.Template, msg.Data)
		}

		// If no template, use data["body"] as raw content
		body := ""
		if msg.Data != nil {
			body = msg.Data["body"]
		}
		return w.emailService.Send(msg.To, msg.Subject, body)
	})
}
