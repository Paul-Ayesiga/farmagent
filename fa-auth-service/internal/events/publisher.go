package events

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EmailQueue = "email_queue"
	SMSQueue   = "sms_queue"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type EmailEvent struct {
	To       string            `json:"to"`
	Subject  string            `json:"subject"`
	Template string            `json:"template"`
	Data     map[string]string `json:"data"`
}

type SMSEvent struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

func NewPublisher(url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare queues
	queues := []string{EmailQueue, SMSQueue}
	for _, q := range queues {
		_, err := ch.QueueDeclare(q, true, false, false, false, nil)
		if err != nil {
			return nil, err
		}
	}

	log.Println("✅ Connected to RabbitMQ (Event Publisher)")
	return &Publisher{conn: conn, channel: ch}, nil
}

func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Publisher) PublishEmail(ctx context.Context, event EmailEvent) error {
	return p.publish(ctx, EmailQueue, event)
}

func (p *Publisher) PublishSMS(ctx context.Context, event SMSEvent) error {
	return p.publish(ctx, SMSQueue, event)
}

func (p *Publisher) publish(ctx context.Context, queue string, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = p.channel.PublishWithContext(
		ctx,
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		log.Printf("Failed to publish to %s: %v", queue, err)
		return err
	}

	log.Printf("📤 Published event to %s", queue)
	return nil
}
