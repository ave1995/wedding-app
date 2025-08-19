package store

import (
	"context"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type AttemptAnswerStore interface {
	GetAnsweredBySession(ctx context.Context, sessionID uuid.UUID) ([]*model.AttemptAnswer, error)
	CreateAttemptAnswer(ctx context.Context, params model.CreateAttemptAnswerParams) (*model.AttemptAnswer, error)
}
