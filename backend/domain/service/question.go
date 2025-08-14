package service

import (
	"context"
	"wedding-app/domain/model"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, text, quizID string) (*model.Question, error)
	GetQuestionByID(ctx context.Context, id string) (*model.Question, error)
	GetQuestionsByQuizID(ctx context.Context, quizID string) ([]*model.Question, error)
}
