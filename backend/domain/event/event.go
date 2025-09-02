package event

import (
	"time"

	"github.com/google/uuid"
)

type AnswerSubmittedEvent struct {
	SessionBaseEvent
	QuestionText string
	Answers      []string
	SubmittedAt  time.Time
}

type SessionBaseEvent struct {
	SessionID   uuid.UUID
	Username    string
	UserIconUrl string
	Quizname    string
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
