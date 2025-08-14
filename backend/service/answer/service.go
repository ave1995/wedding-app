package answer

import (
	"context"
	"fmt"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"

	"github.com/google/uuid"
)

type anwserService struct {
	answerStore store.AnswerStore
}

func NewAnswerServoce(answerStore store.AnswerStore) service.AnswerService {
	return &anwserService{answerStore: answerStore}
}

// CreateAnswer implements service.AnswerService.
func (a *anwserService) CreateAnswer(ctx context.Context, text string, questionID string, isCorrect bool) (*model.Answer, error) {
	parsed, err := uuid.Parse(questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse question ID %q: %w", questionID, err)
	}
	return a.answerStore.CreateAnswer(ctx, text, parsed, isCorrect)
}

// GetAnswerByID implements service.AnswerService.
func (a *anwserService) GetAnswerByID(ctx context.Context, id string) (*model.Answer, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID %q: %w", id, err)
	}
	return a.answerStore.GetAnswerByID(ctx, parsed)
}

// GetAnswersByQuestionID implements service.AnswerService.
func (a *anwserService) GetAnswersByQuestionID(ctx context.Context, questionID string) ([]*model.Answer, error) {
	parsed, err := uuid.Parse(questionID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse question ID %q: %w", questionID, err)
	}
	return a.answerStore.GetAnswersByQuestionID(ctx, parsed)
}
