package event

import (
	"time"

	"github.com/google/uuid"
)

// RADEK: čím si event modely zasloužily vlastní package a nejsou v modelu jako všechny ostaní modely?

type SessionBaseEvent struct {
	SessionID   uuid.UUID
	Username    string
	UserID      uuid.UUID
	UserIconUrl string
	Quizname    string
}

type QuestionOpenedEvent struct {
	SessionBaseEvent
	QuestionID   uuid.UUID
	QuestionText string
	OpenedAt     time.Time
}

type AnswerSubmittedEvent struct {
	SessionBaseEvent
	QuestionID   uuid.UUID
	QuestionText string
	Answers      []string
	SubmittedAt  time.Time
}

// SessionStartEvent embeds SessionBase and adds StartedAt
type SessionStartEvent struct {
	SessionBaseEvent
	StartedAt time.Time
}

// SessionEndEvent embeds SessionBase and adds EndedAt
type SessionEndEvent struct {
	SessionBaseEvent
	EndedAt time.Time
}
