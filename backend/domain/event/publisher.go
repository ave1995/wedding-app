package event

// RADEK: nesedí ti to na struct, co jsme si definovali, toto je pravděpodobně connector
type EventPublisher interface {
	PublishQuestionOpened(e *QuestionOpenedEvent) error
	PublishAnswerSubmitted(e *AnswerSubmittedEvent) error
	PublishSessionStarted(e *SessionStartEvent) error
	PublishSessionEnded(e *SessionEndEvent) error
}
