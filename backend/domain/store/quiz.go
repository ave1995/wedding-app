package store

import (
	"context"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type QuizStore interface {
	CreateQuiz(ctx context.Context, name string) (*model.Quiz, error)
	GetQuizByInviteCode(ctx context.Context, inviteCode uuid.UUID) (*model.Quiz, error)
	GetQuizByID(ctx context.Context, id uuid.UUID) (*model.Quiz, error)
}
