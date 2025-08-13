package model

import "github.com/google/uuid"

type Question struct {
	ID     uuid.UUID
	QuizID uuid.UUID
	Text   string
}
