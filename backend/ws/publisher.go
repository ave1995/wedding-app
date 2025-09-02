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

const TopicAnswerSubmittedEvent = "answer_submit"

// PublishAnswerSubmitted implements event.EventPublisher.
func (p *publisher) PublishAnswerSubmitted(e *event.AnswerSubmittedEvent) error {
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

const TopicSessionStartEvent = "session_start"

func (p *publisher) PublishSessionStarted(e *event.SessionStartEvent) error {
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

const TopicSessionEndEvent = "session_end"

func (p *publisher) PublishSessionEnded(e *event.SessionEndEvent) error {
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
