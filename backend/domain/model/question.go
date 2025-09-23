package model

import (
	"time"

	"github.com/google/uuid"
)

// RADEK: i type aliasy si zaslouží sepre soubor
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
	PhotoPath *string
	Answers   []*Answer
	CreatedAt time.Time
}
