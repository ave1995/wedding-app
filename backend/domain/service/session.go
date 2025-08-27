package service

import (
	"context"
	"wedding-app/domain/model"
)

type SessionService interface {
	StartSession(ctx context.Context, userID string, quizID string) (*model.Session, error)
	GetCurrentQuestion(ctx context.Context, sessionID string) (*model.Question, error)
	SubmitAnswer(ctx context.Context, sessionID string, questionID string, answerID string) (isCompleted bool, err error)
	GetResult(ctx context.Context, sessionID string) (*model.Result, error)
	GetSessionByID(ctx context.Context, sessionID string) (*model.Session, error)
}
