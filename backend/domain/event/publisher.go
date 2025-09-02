package event

type EventPublisher interface {
	PublishAnswerSubmitted(e *AnswerSubmittedEvent) error
	PublishSessionStarted(e *SessionStartEvent) error
	PublishSessionEnded(e *SessionEndEvent) error
}
