package store

import (
	"context"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type QuestionStore interface {
	CreateQuestion(ctx context.Context, text string, quizID uuid.UUID) (*model.Question, error)
	GetQuestionByID(ctx context.Context, id uuid.UUID) (*model.Question, error)
	GetQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]*model.Question, error)
}
