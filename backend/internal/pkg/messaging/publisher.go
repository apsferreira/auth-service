package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/google/uuid"
)

const (
	exchangeName = "auth.events"
	exchangeType = "topic"

	RoutingKeyOTPRequested         = "otp.requested"
	RoutingKeyOTPVerified          = "otp.verified"
	RoutingKeyCustomerRegistered   = "customer.registered"

	maxRetries = 5
)

// OTPRequestedEvent is the payload published when a user requests an OTP code.
// Exchange: auth.events | Routing key: otp.requested
type OTPRequestedEvent struct {
	Email       string `json:"email"`
	ServiceName string `json:"service_name"`
	Channel     string `json:"channel"` // "email" | "telegram"
}

// CustomerRegisteredEvent is published after a successful login/register
// so the customer-service can upsert the customer record.
// Exchange: auth.events | Routing key: customer.registered
type CustomerRegisteredEvent struct {
	AuthUserID string `json:"auth_user_id"`
	Email      string `json:"email"`
	FullName   string `json:"full_name"`
	Username   string `json:"username,omitempty"`
	IsNew      bool   `json:"is_new"`
	OccurredAt string `json:"occurred_at"`
}

// Publisher handles publishing events to the auth.events exchange.
// It is safe for concurrent use.
type Publisher struct {
	rawURL  string
	mu      sync.Mutex
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewPublisher(rawURL string) (*Publisher, error) {
	p := &Publisher{rawURL: rawURL}
	if err := p.connect(); err != nil {
		return nil, err
	}
	return p, nil
}

func sanitizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "<invalid-url>"
	}
	if u.User != nil {
		u.User = url.User("***")
	}
	return u.String()
}

func (p *Publisher) connect() error {
	var err error
	for i := 0; i < maxRetries; i++ {
		p.conn, err = amqp.Dial(p.rawURL)
		if err == nil {
			break
		}
		log.Printf("[auth.events] connection attempt %d/%d failed (url: %s): %v",
			i+1, maxRetries, sanitizeURL(p.rawURL), err)
		time.Sleep(time.Duration(i+1) * 2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ after %d retries: %w", maxRetries, err)
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}

	return p.channel.ExchangeDeclare(
		exchangeName, exchangeType,
		true, false, false, false, nil,
	)
}

// PublishOTPRequested publishes an otp.requested event to the auth.events exchange.
// notification-service consumes this and generates + sends the OTP.
func (p *Publisher) PublishOTPRequested(email, serviceName, channel string) error {
	event := OTPRequestedEvent{
		Email:       email,
		ServiceName: serviceName,
		Channel:     channel,
	}
	return p.publish(RoutingKeyOTPRequested, event)
}

// PublishCustomerRegistered publishes a customer.registered event to the auth.events exchange.
// customer-service consumes this to upsert the customer record after a successful login/register.
func (p *Publisher) PublishCustomerRegistered(evt CustomerRegisteredEvent) error {
	return p.publish(RoutingKeyCustomerRegistered, evt)
}

func (p *Publisher) publish(routingKey string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	msg := amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		MessageId:    uuid.NewString(),
		Timestamp:    time.Now().UTC(),
		Body:         body,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(ctx, exchangeName, routingKey, false, false, msg)
	if err != nil {
		// Try reconnecting once on failure
		log.Printf("[auth.events] publish failed, attempting reconnect: %v", err)
		if reconnErr := p.connect(); reconnErr != nil {
			return fmt.Errorf("publish failed and reconnect failed: %w", err)
		}
		// Fresh context for retry
		retryCtx, retryCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer retryCancel()
		msg.MessageId = uuid.NewString() // new ID for the retry
		return p.channel.PublishWithContext(retryCtx, exchangeName, routingKey, false, false, msg)
	}
	return nil
}

func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
