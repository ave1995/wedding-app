package model

import "github.com/google/uuid"

// Session which User use to complete Quiz
type Session struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	QuizID        uuid.UUID
	CurrentQIndex int
	TotalQCount   int
	IsCompleted   bool
}
