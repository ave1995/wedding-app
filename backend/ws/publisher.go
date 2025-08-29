package ws

import (
	"encoding/json"
	"wedding-app/domain/event"
)

type publisher struct {
	hub *Hub
}

func NewPublisher(h *Hub) event.EventPublisher {
	return &publisher{
		hub: h,
	}
}

const TopicAnswerSubmittedEvent = "answers"

// PublishAnswerSubmitted implements event.EventPublisher.
func (p *publisher) PublishAnswerSubmitted(e event.AnswerSubmittedEvent) error {
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	p.hub.broadcast <- broadcastMessage{
		Topic: TopicAnswerSubmittedEvent,
		Data:  data,
	}

	return nil
}
