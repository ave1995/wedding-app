package session

import (
	"context"
	"fmt"
	"wedding-app/domain/apperrors"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"

	"github.com/google/uuid"
)

type sessionService struct {
	sessionStore       store.SessionStore
	questionStore      store.QuestionStore
	attemptAnswerStore store.AttemptStore
	answerStore        store.AnswerStore
}

func NewSessionService(ss store.SessionStore, qs store.QuestionStore, aas store.AttemptStore, as store.AnswerStore) service.SessionService {
	return &sessionService{sessionStore: ss, questionStore: qs, attemptAnswerStore: aas, answerStore: as}
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
	// Grab the current question by index
	return questions[session.CurrentQIndex], nil
}

// StartQuiz implements service.SessionService.
func (s *sessionService) StartQuiz(ctx context.Context, userID string, quizID string) (*model.Session, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID %q: %w", quizID, err)
	}
	parsedQuizID, err := uuid.Parse(quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quiz ID %q: %w", quizID, err)
	}
	activeSession, err := s.sessionStore.FindActive(ctx, parsedUserID, parsedQuizID)
	if err != nil {
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
	answerID string,
) (isCompleted bool, err error) {
	parsedSessionID, err := uuid.Parse(sessionID)
	if err != nil {
		return false, fmt.Errorf("failed to parse session ID: %w", err)
	}
	parsedQuestionID, err := uuid.Parse(questionID)
	if err != nil {
		return false, fmt.Errorf("failed to parse question ID: %w", err)
	}
	parsedAnswerID, err := uuid.Parse(answerID)
	if err != nil {
		return false, fmt.Errorf("failed to parse answer ID: %w", err)
	}
	// Get session and check
	session, err := s.sessionStore.FindByID(ctx, parsedSessionID)
	if err != nil {
		return false, fmt.Errorf("failed to load session: %w", err)
	}
	if session.IsCompleted {
		return false, fmt.Errorf("session already completed")
	}
	// Get answer for attempt
	answer, err := s.answerStore.GetAnswerByIDAndQuestionID(ctx, parsedAnswerID, parsedQuestionID)
	if err != nil {
		return false, fmt.Errorf("failed to load answer: %w", err)
	}
	// Check if It's already answered
	answered, err := s.attemptAnswerStore.GetAnsweredBySessionIDAndQuestionID(ctx, session.ID, parsedQuestionID)
	if answered != nil {
		return false, fmt.Errorf("question is already answered")
	}
	if err != nil && err != apperrors.ErrNotFound {
		return false, err
	}
	// Creating attempt
	if _, err := s.attemptAnswerStore.CreateAttemptAnswer(ctx, model.CreateAttemptParams{
		SessionID:  session.ID,
		QuestionID: parsedQuestionID,
		AnswerID:   answer.ID,
		IsCorrect:  answer.IsCorrect,
	}); err != nil {
		return false, fmt.Errorf("failed to save attempt answer: %w", err)
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
