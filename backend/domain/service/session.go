package service

import (
	"context"

	"wedding-app/domain/model"
	"wedding-app/dto"
)

type SessionService interface {
	StartSession(ctx context.Context, userID string, quizID string) (*model.Session, error)
	GetCurrentQuestion(ctx context.Context, sessionID string) (*dto.QuestionResponse, error)
	// RADEK: Pojmenované return typy jsou super ale moc se nepoužívají, protože to může způsobovat problémy,
	// Ani jsem nevěděl, že to jde použív u definice interface :D, lepší je nepojmenovat a napsat hezký kommentář/dokumentaci :)
	SubmitAnswer(ctx context.Context, sessionID string, questionID string, answerIDs []string) (isCompleted bool, err error)
	GetResult(ctx context.Context, sessionID string) (*model.Result, error)
	GetSessionByID(ctx context.Context, sessionID string) (*model.Session, error)
	GetActiveSessionsByQuizID(ctx context.Context, quizID string) ([]*model.Session, error)
	GetResultsByQuizID(ctx context.Context, quizID string) ([]*dto.UserResult, error)
}
