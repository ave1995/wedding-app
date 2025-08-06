package quiz

import (
	"context"
	"fmt"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"

	"github.com/google/uuid"
)

type quizService struct {
	store store.QuizStore
}

func NewQuizService(store store.QuizStore) service.QuizService {
	return &quizService{store: store}
}

// CreateQuiz implements service.QuizService.
func (q *quizService) CreateQuiz(ctx context.Context, name string) (*model.Quiz, error) {
	return q.store.CreateQuiz(ctx, name)
}

// GetQuizByInviteCode implements service.QuizService.
func (q *quizService) GetQuizByInviteCode(ctx context.Context, inviteCode string) (*model.Quiz, error) {
	if err := uuid.Validate(inviteCode); err != nil {
		return nil, fmt.Errorf("invalid invite code %q: %w", inviteCode, err)
	}

	return q.store.GetQuizByInviteCode(ctx, inviteCode)
}

// GetQuizByID implements service.QuizService.
func (q *quizService) GetQuizByID(ctx context.Context, id string) (*model.Quiz, error) {
	if err := uuid.Validate(id); err != nil {
		return nil, fmt.Errorf("invalid invite code %q: %w", id, err)
	}

	return q.store.GetQuizByID(ctx, id)
}
