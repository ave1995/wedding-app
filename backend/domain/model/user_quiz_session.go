package model

import "github.com/google/uuid"

type UserQuizSession struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	QuizID      uuid.UUID
	CurrentQID  uuid.UUID
	IsCompleted bool
}
