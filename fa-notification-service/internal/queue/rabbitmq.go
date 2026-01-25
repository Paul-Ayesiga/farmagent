package queue

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EmailQueue = "email_queue"
	SMSQueue   = "sms_queue"
	PushQueue  = "push_queue"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type EmailMessage struct {
	To       string            `json:"to"`
	Subject  string            `json:"subject"`
	Template string            `json:"template"`
	Data     map[string]string `json:"data"`
}

type SMSMessage struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

type PushMessage struct {
	DeviceToken string            `json:"device_token"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Data        map[string]string `json:"data,omitempty"`
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
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
	queues := []string{EmailQueue, SMSQueue, PushQueue}
	for _, q := range queues {
		_, err := ch.QueueDeclare(
			q,     // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return nil, err
		}
	}

	log.Println("✅ Connected to RabbitMQ")
	return &RabbitMQ{conn: conn, channel: ch}, nil
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

// Publish methods

func (r *RabbitMQ) PublishEmail(ctx context.Context, msg EmailMessage) error {
	return r.publish(ctx, EmailQueue, msg)
}

func (r *RabbitMQ) PublishSMS(ctx context.Context, msg SMSMessage) error {
	return r.publish(ctx, SMSQueue, msg)
}

func (r *RabbitMQ) PublishPush(ctx context.Context, msg PushMessage) error {
	return r.publish(ctx, PushQueue, msg)
}

func (r *RabbitMQ) publish(ctx context.Context, queue string, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(
		ctx,
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

// Consume methods

func (r *RabbitMQ) ConsumeEmails(handler func(EmailMessage) error) error {
	return r.consume(EmailQueue, func(body []byte) error {
		var msg EmailMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return err
		}
		return handler(msg)
	})
}

func (r *RabbitMQ) ConsumeSMS(handler func(SMSMessage) error) error {
	return r.consume(SMSQueue, func(body []byte) error {
		var msg SMSMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return err
		}
		return handler(msg)
	})
}

func (r *RabbitMQ) ConsumePush(handler func(PushMessage) error) error {
	return r.consume(PushQueue, func(body []byte) error {
		var msg PushMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return err
		}
		return handler(msg)
	})
}

func (r *RabbitMQ) consume(queue string, handler func([]byte) error) error {
	msgs, err := r.channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error processing message: %v", err)
				d.Nack(false, true) // requeue
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}
