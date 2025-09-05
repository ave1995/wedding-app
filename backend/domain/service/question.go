package service

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/dto"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, text, quizID string, questionType model.QuestionType, photoPath *string) (*model.Question, error)
	GetQuestionByID(ctx context.Context, id string) (*model.Question, error)
	GetQuestionsByQuizID(ctx context.Context, quizID string) ([]*model.Question, error)
	RevealQuestionByQuizIDAndIndex(ctx context.Context, quizID string, index int) (*dto.RevealResponse, error)
	RevealQuestionStatsByID(ctx context.Context, id string) (*dto.StatsResponse, error)
}
