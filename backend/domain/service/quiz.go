package service

import (
	"context"
	"wedding-app/domain/model"
)

type QuizService interface {
	CreateQuiz(ctx context.Context, name string) (*model.Quiz, error)
	GetQuizByInviteCode(ctx context.Context, inviteCode string) (*model.Quiz, error)
	GetQuizByID(ctx context.Context, id string) (*model.Quiz, error)
}
