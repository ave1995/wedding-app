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

type AnswerService interface {
	CreateAnswer(ctx context.Context, text, questionID string) (*model.Answer, error)
	GetAnswerByID(ctx context.Context, id string) (*model.Answer, error)
	GetAnswersByQuestionID(ctx context.Context, questionID string) ([]*model.Answer, error)
}
