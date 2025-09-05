package question

import (
	"context"
	"fmt"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/dto"
	"wedding-app/utils"

	"github.com/google/uuid"
)

type questionService struct {
	questionStore      store.QuestionStore
	sessionStore       store.SessionStore
	answerStore        store.AnswerStore
	attemptAnswerStore store.AttemptStore
}

func NewQuestionService(questionStore store.QuestionStore, sessionStore store.SessionStore, answerStore store.AnswerStore, attemptAnswerStore store.AttemptStore) service.QuestionService {
	return &questionService{questionStore: questionStore, sessionStore: sessionStore, answerStore: answerStore, attemptAnswerStore: attemptAnswerStore}
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

func (q *questionService) RevealQuestionStatsByID(ctx context.Context, id string) (*dto.StatsResponse, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse question ID %q: %w", id, err)
	}

	question, err := q.questionStore.GetQuestionByID(ctx, parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to load question %w", err)
	}

	sessions, err := q.sessionStore.GetSessionsByQuizID(ctx, question.QuizID)
	if err != nil {
		return nil, fmt.Errorf("failed to load sessions %w", err)
	}

	correctAnswers, err := q.answerStore.GetCorrectAnswersByQuestionID(ctx, question.ID) // TODO: měl bych využít MongoDB a i cache
	if err != nil {
		return nil, fmt.Errorf("failed to get correct answers: %w", err)
	}

	right := 0
	wrong := 0
	for _, s := range sessions {
		if !s.IsCompleted {
			continue
		}
		answered, err := q.attemptAnswerStore.GetAnsweredBySessionIDAndQuestionID(ctx, s.ID, question.ID)
		if err != nil {
			return nil, err
		}

		correct, err := utils.IsQuestionCorrect(question, answered, correctAnswers)
		if err != nil {
			return nil, err
		}
		if correct {
			right++
		} else {
			wrong++
		}
	}

	return &dto.StatsResponse{Right: right, Wrong: wrong}, nil

}
