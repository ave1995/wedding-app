package question

import (
	"context"
	"fmt"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/dto"

	"github.com/google/uuid"
)

type questionService struct {
	questionStore store.QuestionStore
	answerStore   store.AnswerStore
}

func NewQuestionService(questionStore store.QuestionStore, asnwerStore store.AnswerStore) service.QuestionService {
	return &questionService{questionStore: questionStore, answerStore: asnwerStore}
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

// RevealQuestionByQuizIDAndIndex implements service.QuestionService.
func (q *questionService) RevealQuestionByQuizIDAndIndex(ctx context.Context, quizID string, index int) (*dto.RevealResponse, error) {
	parsed, err := uuid.Parse(quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quiz ID %q: %w", quizID, err)
	}

	questions, err := q.questionStore.GetOrderedQuestionsByQuizID(ctx, parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to load questions: %w", err)
	}

	lenQuestions := len(questions)

	if lenQuestions == 0 {
		return nil, fmt.Errorf("no questions found for quiz ID %s", parsed)
	}
	if index < 0 || index >= lenQuestions {
		return nil, fmt.Errorf("current question index %d out of range (0-%d)", index, lenQuestions-1)
	}

	question := questions[index]

	answers, err := q.answerStore.GetAnswersByQuestionID(ctx, question.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load answers for question %q: %w", question.Text, err)
	}

	question.Answers = answers

	reveal := &dto.RevealResponse{
		Question:    question,
		GoNext:      index < lenQuestions,
		NextIndex:   index + 1,
		TotalQCount: lenQuestions,
	}

	return reveal, nil
}
