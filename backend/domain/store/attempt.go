package store

import (
	"context"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type AttemptStore interface {
	GetAnsweredBySessionIDAndQuestionID(ctx context.Context, sessionID uuid.UUID, questionID uuid.UUID) ([]*model.Attempt, error)
	GetAnsweredBySessionID(ctx context.Context, sessionID uuid.UUID) ([]*model.Attempt, error)
	CreateAttemptAnswer(ctx context.Context, params model.CreateAttemptParams) (*model.Attempt, error)
}
