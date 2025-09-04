package question

import (
	"context"
	"fmt"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"

	"github.com/google/uuid"
)

type questionService struct {
	questionStore store.QuestionStore
}

func NewQuestionService(questionStore store.QuestionStore) service.QuestionService {
	return &questionService{questionStore: questionStore}
}

// CreateQuestion implements service.QuestionService.
func (q *questionService) CreateQuestion(ctx context.Context, text string, quizID string, questionType model.QuestionType, photoPath *string) (*model.Question, error) {
	parsed, err := uuid.Parse(quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quiz ID %q: %w", quizID, err)
	}
	return q.questionStore.CreateQuestion(ctx, text, parsed, questionType, photoPath)
}

// GetQuestionByID implements service.QuestionService.
func (q *questionService) GetQuestionByID(ctx context.Context, id string) (*model.Question, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID %q: %w", id, err)
	}
	return q.questionStore.GetQuestionByID(ctx, parsed)
}

// GetQuestionsByQuizID implements service.QuestionService.
func (q *questionService) GetQuestionsByQuizID(ctx context.Context, quizID string) ([]*model.Question, error) {
	parsed, err := uuid.Parse(quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quiz ID %q: %w", quizID, err)
	}
	return q.questionStore.GetQuestionsByQuizID(ctx, parsed)
}
