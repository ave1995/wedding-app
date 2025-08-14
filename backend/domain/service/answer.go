package service

import (
	"context"
	"wedding-app/domain/model"
)

type AnswerService interface {
	CreateAnswer(ctx context.Context, text, questionID string, isCorrect bool) (*model.Answer, error)
	GetAnswerByID(ctx context.Context, id string) (*model.Answer, error)
	GetAnswersByQuestionID(ctx context.Context, questionID string) ([]*model.Answer, error)
}
