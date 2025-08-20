package model

import "github.com/google/uuid"

// Attempt in Sesssion for Question with selected Answer
type Attempt struct {
	ID         uuid.UUID
	SessionID  uuid.UUID
	QuestionID uuid.UUID
	AnswerID   uuid.UUID
	IsCorrect  bool
}

type CreateAttemptParams struct {
	SessionID  uuid.UUID
	QuestionID uuid.UUID
	AnswerID   uuid.UUID
	IsCorrect  bool
}
