package messaging

import (
	"fmt"
)

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

type MessagingService interface {
	PublishMessage(messageEvent MessageEvent) error
	Close() error
}
