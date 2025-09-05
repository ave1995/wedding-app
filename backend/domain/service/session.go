package service

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/dto"
)

type SessionService interface {
	StartSession(ctx context.Context, userID string, quizID string) (*model.Session, error)
	GetCurrentQuestion(ctx context.Context, sessionID string) (*dto.QuestionResponse, error)
	SubmitAnswer(ctx context.Context, sessionID string, questionID string, answerIDs []string) (isCompleted bool, err error)
	GetResult(ctx context.Context, sessionID string) (*model.Result, error)
	GetSessionByID(ctx context.Context, sessionID string) (*model.Session, error)
	GetActiveSessionsByQuizID(ctx context.Context, quizID string) ([]*model.Session, error)
	GetResultsByQuizID(ctx context.Context, quizID string) ([]*dto.UserResult, error)
}
