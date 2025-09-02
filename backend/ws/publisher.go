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

type envelope[T any] struct {
	Topic string `json:"topic"`
	Data  T      `json:"data"`
}

func wrapEvent[T any](topic string, e T) ([]byte, error) {
	return json.Marshal(envelope[T]{
		Topic: topic,
		Data:  e,
	})
}

const TopicAnswerSubmittedEvent = "answer_submit"

// PublishAnswerSubmitted implements event.EventPublisher.
func (p *publisher) PublishAnswerSubmitted(e *event.AnswerSubmittedEvent) error {
	data, err := wrapEvent(TopicAnswerSubmittedEvent, e)
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
	data, err := wrapEvent(TopicSessionStartEvent, e)
	if err != nil {
		return err
	}

	p.hub.broadcast <- broadcastMessage{
		Topic: TopicSessionStartEvent,
		Data:  data,
	}

	return nil
}

const TopicSessionEndEvent = "session_end"

func (p *publisher) PublishSessionEnded(e *event.SessionEndEvent) error {
	data, err := wrapEvent(TopicSessionEndEvent, e)
	if err != nil {
		return err
	}

	p.hub.broadcast <- broadcastMessage{
		Topic: TopicSessionEndEvent,
		Data:  data,
	}

	return nil
}
