package messaging

import (
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQService struct {
	exchange string
	conn     *amqp.Connection
	ch       *amqp.Channel
}

func NewRabbitMQService(exchange string) (MessagingService, error) {
	var err error

	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return nil, fmt.Errorf("RABBITMQ_URL not set")
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &rabbitMQService{
		exchange: exchange,
		conn:     conn,
		ch:       ch,
	}, nil
}

func (s *rabbitMQService) PublishMessage(msg MessageEvent) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	routingKey := msg.RoutingKey()

	return s.ch.Publish(
		s.exchange,
		routingKey,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *rabbitMQService) Close() error {
	if err := r.ch.Close(); err != nil {
		return err
	}

	if err := r.conn.Close(); err != nil {
		return err
	}
	
	return nil
}
