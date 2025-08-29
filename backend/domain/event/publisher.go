package event

type EventPublisher interface {
	PublishAnswerSubmitted(e AnswerSubmittedEvent) error
}
