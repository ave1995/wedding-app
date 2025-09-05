package store

import (
	"context"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type SessionStore interface {
	CreateSession(ctx context.Context, userID uuid.UUID, quizID uuid.UUID, questionCount int) (*model.Session, error)
	FindActive(ctx context.Context, userID uuid.UUID, quizID uuid.UUID) (*model.Session, error)
	FindByID(ctx context.Context, sessionID uuid.UUID) (*model.Session, error)
	UpdateSession(ctx context.Context, session *model.Session) error
	GetSessionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]*model.Session, error)
}
