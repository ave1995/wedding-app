package session

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
)

type sessionService struct {
}

func NewSessionService() service.SessionService {
	return &sessionService{}
}

// GetCurrentQuestion implements service.SessionService.
func (s *sessionService) GetCurrentQuestion(ctx context.Context, sessionID string) (*model.Question, error) {
	panic("unimplemented")
}

// StartQuiz implements service.SessionService.
func (s *sessionService) StartQuiz(ctx context.Context, userID string, quizID string) (*model.UserQuizSession, error) {
	panic("unimplemented")
}

// SubmitAnswer implements service.SessionService.
func (s *sessionService) SubmitAnswer(ctx context.Context, sessionID string, questionID string, answerID string) (isCompleted bool, err error) {
	panic("unimplemented")
}
