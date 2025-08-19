package store

import (
	"context"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type SessionStore interface {
	CreateSession(ctx context.Context, userID uuid.UUID, quizID uuid.UUID) (*model.UserQuizSession, error)
	FindActive(ctx context.Context, userID uuid.UUID, quizID uuid.UUID) (*model.UserQuizSession, error)
	FindByID(ctx context.Context, sessionID uuid.UUID) (*model.UserQuizSession, error)
	UpdateSession(ctx context.Context, session *model.UserQuizSession) error
}
