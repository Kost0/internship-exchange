package publisher

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func New(url string) (*Publisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	for _, queue := range []string{"application.created", "application.status_changed"} {
		_, err = ch.QueueDeclare(queue, true, false, false, false, nil)
		if err != nil {
			return nil, err
		}
	}

	return &Publisher{conn: conn, channel: ch}, nil
}

func (p *Publisher) Close() {
	p.channel.Close()
	p.conn.Close()
}

func (p *Publisher) PublishApplicationCreated(ctx context.Context, applicationID, studentID, listingID, companyID string) {
	p.publish(ctx, "application.created", map[string]any{
		"applicationId": applicationID,
		"studentId":     studentID,
		"listingId":     listingID,
		"companyId":     companyID,
	})
}

func (p *Publisher) PublishStatusChanged(ctx context.Context, applicationID, studentID, listingID, oldStatus, newStatus, comment, email string) {
	p.publish(ctx, "application.status_changed", map[string]any{
		"applicationId": applicationID,
		"studentId":     studentID,
		"listingId":     listingID,
		"oldStatus":     oldStatus,
		"newStatus":     newStatus,
		"comment":       comment,
		"email":         email,
	})
}

func (p *Publisher) publish(ctx context.Context, queue string, payload map[string]any) {
	event := Event{Type: queue, Payload: payload}
	body, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal event: %v", err)

		return
	}

	err = p.channel.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
	})
	if err != nil {
		log.Printf("failed to publish event to %s: %v", queue, err)
	}
}
