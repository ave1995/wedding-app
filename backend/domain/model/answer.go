package model

import "github.com/google/uuid"

type Answer struct {
	ID         uuid.UUID
	QuestionID uuid.UUID
	Text       string
	IsCorrect  bool
}
