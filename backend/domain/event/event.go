package event

import (
	"time"

	"github.com/google/uuid"
)

type AnswerSubmittedEvent struct {
	SessionID   uuid.UUID
	UserID      uuid.UUID
	QuizID      uuid.UUID
	QuestionID  uuid.UUID
	AnswerIDs   []uuid.UUID
	SubmittedAt time.Time
}
