package store

import (
	"context"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type AnswerStore interface {
	CreateAnswer(ctx context.Context, text string, questionID uuid.UUID, isCorrect bool) (*model.Answer, error)
	GetAnswerByID(ctx context.Context, id uuid.UUID) (*model.Answer, error)
	GetAnswerByIDAndQuestionID(ctx context.Context, answerID uuid.UUID, questionID uuid.UUID) (*model.Answer, error)
	GetAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*model.Answer, error)
	GetCorrectAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*model.Answer, error)
}
