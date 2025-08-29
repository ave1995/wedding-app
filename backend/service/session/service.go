package session

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"wedding-app/domain/apperrors"
	"wedding-app/domain/event"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/utils"

	"github.com/google/uuid"
)

type sessionService struct {
	sessionStore       store.SessionStore
	questionStore      store.QuestionStore
	attemptAnswerStore store.AttemptStore
	answerStore        store.AnswerStore
	publisher          event.EventPublisher
	logger             *slog.Logger
}

func NewSessionService(ss store.SessionStore, qs store.QuestionStore, aas store.AttemptStore, as store.AnswerStore, pub event.EventPublisher, logger *slog.Logger) service.SessionService {
	return &sessionService{sessionStore: ss, questionStore: qs, attemptAnswerStore: aas, answerStore: as, publisher: pub, logger: logger}
}

// GetCurrentQuestion implements service.SessionService.
func (s *sessionService) GetCurrentQuestion(ctx context.Context, sessionID string) (*model.Question, error) {
	parsed, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session ID %q: %w", sessionID, err)
	}
	// Load Quiz Session
	session, err := s.sessionStore.FindByID(ctx, parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to load session: %w", err)
	}
	// Get All Questions for Quiz
	questions, err := s.questionStore.GetOrderedQuestionsByQuizID(ctx, session.QuizID)
	if err != nil {
		return nil, fmt.Errorf("failed to load questions: %w", err)
	}
	if len(questions) == 0 {
		return nil, fmt.Errorf("no questions found for quiz ID %s", session.QuizID)
	}
	if session.CurrentQIndex < 0 || session.CurrentQIndex >= len(questions) {
		return nil, fmt.Errorf("current question index %d out of range (0-%d)", session.CurrentQIndex, len(questions)-1)
	}
	// TODO: zeptat se Radka
	originalQuestion := questions[session.CurrentQIndex]
	question := &model.Question{
		ID:        originalQuestion.ID,
		QuizID:    originalQuestion.QuizID,
		Text:      originalQuestion.Text,
		CreatedAt: originalQuestion.CreatedAt,
	}

	answers, err := s.answerStore.GetAnswersByQuestionID(ctx, question.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load answers for question %q: %w", question.Text, err)
	}

	question.Answers = answers
	// Grab the current question by index
	return question, nil
}

// StartSession implements service.SessionService.
func (s *sessionService) StartSession(ctx context.Context, userID string, quizID string) (*model.Session, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID %q: %w", quizID, err)
	}
	parsedQuizID, err := uuid.Parse(quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quiz ID %q: %w", quizID, err)
	}
	activeSession, err := s.sessionStore.FindActive(ctx, parsedUserID, parsedQuizID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, err
	}
	if activeSession != nil {
		return activeSession, nil // return existing active session
	}
	count, err := s.questionStore.GetCountQuestionsByQuizID(ctx, parsedQuizID)
	if err != nil {
		return nil, err
	}
	return s.sessionStore.CreateSession(ctx, parsedUserID, parsedQuizID, count)
}

// SubmitAnswer implements service.SessionService.
func (s *sessionService) SubmitAnswer(
	ctx context.Context,
	sessionID string,
	questionID string,
	answerIDs []string,
) (isCompleted bool, err error) {
	// Parse session ID
	parsedSessionID, err := uuid.Parse(sessionID)
	if err != nil {
		return false, fmt.Errorf("failed to parse session ID: %w", err)
	}

	// Parse question ID
	parsedQuestionID, err := uuid.Parse(questionID)
	if err != nil {
		return false, fmt.Errorf("failed to parse question ID: %w", err)
	}

	// Parse all answer IDs
	parsedAnswerIDs := make([]uuid.UUID, len(answerIDs))
	for i, id := range answerIDs {
		parsedAnswerIDs[i], err = uuid.Parse(id)
		if err != nil {
			return false, fmt.Errorf("invalid answer ID: %s", id)
		}
	}

	// Load session
	session, err := s.sessionStore.FindByID(ctx, parsedSessionID)
	if err != nil {
		return false, fmt.Errorf("failed to load session: %w", err)
	}
	if session.IsCompleted {
		return false, apperrors.ErrSessionCompleted
	}

	// Load question to check type
	question, err := s.questionStore.GetQuestionByID(ctx, parsedQuestionID)
	if err != nil {
		return false, fmt.Errorf("failed to load question: %w", err)
	}

	// Validation: single-choice questions must have exactly one answer
	if question.Type == model.SingleChoice && len(parsedAnswerIDs) != 1 {
		return false, fmt.Errorf("single-choice question must have exactly one selected answer")
	}

	// Check if question was already answered
	answered, err := s.attemptAnswerStore.GetAnsweredBySessionIDAndQuestionID(ctx, session.ID, parsedQuestionID)
	if answered != nil {
		return false, fmt.Errorf("question is already answered")
	}
	if err != nil && err != apperrors.ErrNotFound {
		return false, err
	}

	// Save all selected answers
	for _, ansID := range parsedAnswerIDs {
		answer, err := s.answerStore.GetAnswerByIDAndQuestionID(ctx, ansID, parsedQuestionID)
		if err != nil {
			return false, fmt.Errorf("failed to load answer: %w", err)
		}

		if _, err := s.attemptAnswerStore.CreateAttemptAnswer(ctx, model.CreateAttemptParams{
			SessionID:  session.ID,
			QuestionID: parsedQuestionID,
			AnswerID:   answer.ID,
			IsCorrect:  answer.IsCorrect,
		}); err != nil {
			return false, fmt.Errorf("failed to save attempt answer: %w", err)
		}
	}

	if err := s.publisher.PublishAnswerSubmitted(event.AnswerSubmittedEvent{
		SessionID:   session.ID,
		UserID:      session.UserID,
		QuizID:      session.QuizID,
		QuestionID:  question.ID,
		AnswerIDs:   parsedAnswerIDs,
		SubmittedAt: time.Now(),
	}); err != nil {
		s.logger.Error("failed to publish AnswerSubmittedEvent: %v", utils.ErrAttr(err))
	}

	// Move to next question
	session.CurrentQIndex++
	if session.CurrentQIndex >= session.TotalQCount {
		session.IsCompleted = true
		isCompleted = true
	}

	if err := s.sessionStore.UpdateSession(ctx, session); err != nil {
		return false, fmt.Errorf("failed to update session: %w", err)
	}

	return isCompleted, nil
}

// GetResult implements service.SessionService.
func (s *sessionService) GetResult(ctx context.Context, sessionID string) (*model.Result, error) {
	parsedSessionID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session ID: %w", err)
	}
	session, err := s.sessionStore.FindByID(ctx, parsedSessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	attempts, err := s.attemptAnswerStore.GetAnsweredBySessionID(ctx, parsedSessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attempts: %w", err)
	}

	score := 0
	for _, a := range attempts {
		if a.IsCorrect {
			score++
		}
	}

	total := session.TotalQCount
	percentage := 0
	if total > 0 {
		percentage = score * 100 / total
	}

	result := &model.Result{
		Score:      score,
		Total:      total,
		Percentage: percentage,
	}
	return result, nil
}

// GetSessionByID implements service.SessionService.
func (s *sessionService) GetSessionByID(ctx context.Context, sessionID string) (*model.Session, error) {
	parsedSessionID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session ID: %w", err)
	}
	return s.sessionStore.FindByID(ctx, parsedSessionID)
}
