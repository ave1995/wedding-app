package service

import (
	"context"
	"wedding-app/domain/model"
)

type SessionService interface {
	StartQuiz(ctx context.Context, userID string, quizID string) (*model.UserQuizSession, error)
	GetCurrentQuestion(ctx context.Context, sessionID string) (*model.Question, error)
	SubmitAnswer(ctx context.Context, sessionID string, questionID string, answerID string) (isCompleted bool, err error)
}
