package model

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	Token     string
	ExpiresAt time.Time
	UserID    uuid.UUID
	IsGuest   bool
	QuizID    uuid.UUID
}
