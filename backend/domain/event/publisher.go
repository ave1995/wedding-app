package event

type EventPublisher interface {
	PublishQuestionOpened(e *QuestionOpenedEvent) error
	PublishAnswerSubmitted(e *AnswerSubmittedEvent) error
	PublishSessionStarted(e *SessionStartEvent) error
	PublishSessionEnded(e *SessionEndEvent) error
}
