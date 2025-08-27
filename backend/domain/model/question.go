package model

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID        uuid.UUID
	QuizID    uuid.UUID
	Text      string
	Answers   []*Answer
	CreatedAt time.Time
}
