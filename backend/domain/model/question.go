package model

import "github.com/google/uuid"

type Answer struct {
	Text      string
	IsCorrect bool
}

type Question struct {
	ID      uuid.UUID
	Text    string
	Answers []Answer
}
