package model

import (
	"time"

	"github.com/google/uuid"
)

type QuestionType string

const (
	SingleChoice   QuestionType = "single_choice"
	MultipleChoice QuestionType = "multiple_choice"
)

type Question struct {
	ID        uuid.UUID
	QuizID    uuid.UUID
	Text      string
	Type      QuestionType
	PhotoPath string
	Answers   []*Answer
	CreatedAt time.Time
}
