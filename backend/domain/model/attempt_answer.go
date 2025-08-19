package model

import "github.com/google/uuid"

type AttemptAnswer struct {
	ID         uuid.UUID
	SessionID  uuid.UUID
	QuestionID uuid.UUID
	AnswerID   uuid.UUID
	IsCorrect  bool
}

type CreateAttemptAnswerParams struct {
	SessionID  uuid.UUID
	QuestionID uuid.UUID
	AnswerID   uuid.UUID
	IsCorrect  bool
}
