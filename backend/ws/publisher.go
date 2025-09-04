package ws

import (
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

const TopicQuestionOpenedEvent = "question_open"

// PublishQuestionOpened implements event.EventPublisher.
func (p *publisher) PublishQuestionOpened(e *event.QuestionOpenedEvent) error {
	data, err := wrapEvent(TopicQuestionOpenedEvent, e)
	if err != nil {
		return err
	}

	p.hub.broadcast <- broadcastMessage{
		Topic: TopicQuestionOpenedEvent,
		Data:  data,
	}

	return nil
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
