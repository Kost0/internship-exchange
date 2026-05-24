package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/Kost0/internship-exchange/services/notification-service/internal/mailer"
)

type Event struct {
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	mailer  *mailer.Mailer
}

func New(url string, m *mailer.Mailer) (*Consumer, error) {
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

	return &Consumer{conn: conn, channel: ch, mailer: m}, nil
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.conn.Close()
}

func (c *Consumer) Start(ctx context.Context) error {
	createdMsgs, err := c.channel.Consume("application.created", "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	changedMsgs, err := c.channel.Consume("application.status_changed", "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Println("notification-service consuming events...")

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-createdMsgs:
			if !ok {
				return nil
			}
			c.handleApplicationCreated(msg)
		case msg, ok := <-changedMsgs:
			if !ok {
				return nil
			}
			c.handleStatusChanged(msg)
		}
	}
}

func (c *Consumer) handleApplicationCreated(msg amqp.Delivery) {
	var event Event
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("failed to unmarshal application.created: %v", err)
		msg.Nack(false, false)
		return
	}

	applicationID, _ := event.Payload["applicationId"].(string)
	listingID, _ := event.Payload["listingId"].(string)
	companyEmail, _ := event.Payload["companyEmail"].(string)

	log.Printf("new application %s for listing %s", applicationID, listingID)

	if companyEmail == "" {
		log.Printf("no company email in payload, skipping send")
		msg.Ack(false)
		return
	}

	subject := "Новый отклик на вашу вакансию"
	body := fmt.Sprintf(`
		<h2>Новый отклик</h2>
		<p>На вашу вакансию поступил новый отклик.</p>
		<p>ID отклика: <strong>%s</strong></p>
		<p>Войдите в <a href="http://localhost:3000/dashboard/company">личный кабинет</a> чтобы просмотреть профиль кандидата.</p>
	`, applicationID)

	if err := c.mailer.Send(companyEmail, subject, body); err != nil {
		log.Printf("failed to send email to %s: %v", companyEmail, err)
	} else {
		log.Printf("email sent to %s for application %s", companyEmail, applicationID)
	}

	msg.Ack(false)
}

func (c *Consumer) handleStatusChanged(msg amqp.Delivery) {
	var event Event
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Printf("failed to unmarshal application.status_changed: %v", err)
		msg.Nack(false, false)
		return
	}

	applicationID, _ := event.Payload["applicationId"].(string)
	newStatus, _ := event.Payload["newStatus"].(string)
	comment, _ := event.Payload["comment"].(string)
	studentEmail, _ := event.Payload["studentEmail"].(string)

	log.Printf("application %s status changed to %s", applicationID, newStatus)

	if studentEmail == "" {
		log.Printf("no student email in payload, skipping send")
		msg.Ack(false)
		return
	}

	statusLabels := map[string]string{
		"reviewing": "На рассмотрении",
		"interview": "Приглашение на интервью",
		"accepted":  "Принят",
		"rejected":  "Отказ",
	}

	label, ok := statusLabels[newStatus]
	if !ok {
		msg.Ack(false)
		return
	}

	subject := fmt.Sprintf("Статус вашего отклика изменён: %s", label)
	body := buildStatusEmail(label, comment)

	if err := c.mailer.Send(studentEmail, subject, body); err != nil {
		log.Printf("failed to send email to %s: %v", studentEmail, err)
	} else {
		log.Printf("email sent to %s, status: %s", studentEmail, newStatus)
	}

	msg.Ack(false)
}

func buildStatusEmail(statusLabel, comment string) string {
	commentBlock := ""
	if comment != "" {
		commentBlock = fmt.Sprintf(`<p>Комментарий: <em>%s</em></p>`, comment)
	}

	return fmt.Sprintf(`
		<h2>Статус вашего отклика обновлён</h2>
		<p>Новый статус: <strong>%s</strong></p>
		%s
		<p>Войдите в <a href="http://localhost:3000/dashboard/student">личный кабинет</a> чтобы увидеть детали.</p>
	`, statusLabel, commentBlock)
}
