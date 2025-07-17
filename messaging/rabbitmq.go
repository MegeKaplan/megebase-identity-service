package messaging

import (
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	// q    amqp.Queue
)

func StartProducer() error {
	var err error

	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return fmt.Errorf("RABBITMQ_URL not set")
	}

	conn, err = amqp.Dial(url)
	if err != nil {
		return err
	}

	ch, err = conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		"megebase.topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func CloseProducer() error {
	if ch != nil {
		if err := ch.Close(); err != nil {
			return err
		}
	}
	if conn != nil {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

type MessageEvent struct {
	Service string                 `json:"service"`
	Entity  string                 `json:"entity"`
	Action  string                 `json:"action"`
	Channel string                 `json:"channel,omitempty"`
	To      string                 `json:"to,omitempty"`
	Data    map[string]interface{} `json:"data"`
}

func (m *MessageEvent) RoutingKey() string {
	if m.Channel != "" {
		return fmt.Sprintf("%s.%s.%s.%s", m.Service, m.Entity, m.Action, m.Channel)
	}
	return fmt.Sprintf("%s.%s.%s", m.Service, m.Entity, m.Action)
}

func PublishMessage(msg MessageEvent) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	routingKey := msg.RoutingKey()

	return ch.Publish(
		"megebase.topic",
		routingKey,
		true,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
