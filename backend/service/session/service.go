package session

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"wedding-app/assembler"
	"wedding-app/domain/apperrors"
	"wedding-app/domain/event"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/dto"
	"wedding-app/utils"

	"github.com/google/uuid"
)

type sessionService struct {
	sessionStore       store.SessionStore
	questionStore      store.QuestionStore
	attemptAnswerStore store.AttemptStore
	answerStore        store.AnswerStore
	assembler          *assembler.Assembler
	publisher          event.EventPublisher
	logger             *slog.Logger
}

func NewSessionService(ss store.SessionStore, qs store.QuestionStore, aas store.AttemptStore, as store.AnswerStore, a *assembler.Assembler, pub event.EventPublisher, logger *slog.Logger) service.SessionService {
	return &sessionService{sessionStore: ss, questionStore: qs, attemptAnswerStore: aas, answerStore: as, assembler: a, publisher: pub, logger: logger}
}

// GetCurrentQuestion implements service.SessionService.
func (s *sessionService) GetCurrentQuestion(ctx context.Context, sessionID string) (*dto.QuestionResponse, error) {
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

	question := questions[session.CurrentQIndex]

	answers, err := s.answerStore.GetAnswersByQuestionID(ctx, question.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load answers for question %q: %w", question.Text, err)
	}

	// TODO: udělat DTO
	// Without reveal
	for _, a := range answers {
		a.IsCorrect = false
	}

	question.Answers = answers

	response := &dto.QuestionResponse{
		SessionID:     session.ID,
		Completed:     session.IsCompleted,
		Question:      question,
		CurrentQIndex: session.CurrentQIndex + 1,
		TotalQCount:   len(questions),
	}

	assembledEvent, err := s.assembler.ToQuestionOpenedEvent(ctx, session.ID, session.UserID, session.QuizID, question.ID)
	if err != nil {
		s.logger.Error("failed to assemble QuestionOpenedEvent: %v", utils.ErrAttr(err))
	}
	if err := s.publisher.PublishQuestionOpened(assembledEvent); err != nil {
		s.logger.Error("failed to publish QuestionOpenedEvent: %v", utils.ErrAttr(err))
	}
	// Grab the current question by index
	return response, nil
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
	newSession, err := s.sessionStore.CreateSession(ctx, parsedUserID, parsedQuizID, count)
	if err != nil {
		return nil, err
	}

	assembledEvent, err := s.assembler.ToSessionStartEvent(ctx, newSession.ID, parsedUserID, parsedQuizID)
	if err != nil {
		s.logger.Error("failed to assemble SessionStartEvent: %v", utils.ErrAttr(err))
	}
	if err := s.publisher.PublishSessionStarted(assembledEvent); err != nil {
		s.logger.Error("failed to publish SessionStartEvent: %v", utils.ErrAttr(err))
	}

	return newSession, err
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

	if len(answerIDs) == 0 {
		return false, apperrors.ErrAnswerNotSubmitted
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
	if len(answered) != 0 {
		return false, fmt.Errorf("question is already answered")
	}
	if err != nil {
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

	assembledEvent, err := s.assembler.ToAnswerSubmittedEvent(ctx, session.ID, session.UserID, session.QuizID, question.ID, parsedAnswerIDs)
	if err != nil {
		s.logger.Error("failed to assemble SessionStartEvent: %v", utils.ErrAttr(err))
	}
	if err := s.publisher.PublishAnswerSubmitted(assembledEvent); err != nil {
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

	if isCompleted {
		assembledEvent, err := s.assembler.ToSessionEndEvent(ctx, session.ID, session.UserID, session.QuizID)
		if err != nil {
			s.logger.Error("failed to assemble SessionEndEvent: %v", utils.ErrAttr(err))
		}
		if err := s.publisher.PublishSessionEnded(assembledEvent); err != nil {
			s.logger.Error("failed to publish SessionEndEvent: %v", utils.ErrAttr(err))
		}
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

	allUserAnswers, err := s.attemptAnswerStore.GetAnsweredBySessionID(ctx, parsedSessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attempts: %w", err)
	}

	userAnswersByQuestion := make(map[uuid.UUID][]*model.Attempt)
	for _, a := range allUserAnswers {
		userAnswersByQuestion[a.QuestionID] = append(userAnswersByQuestion[a.QuestionID], a)
	}

	questions, err := s.questionStore.GetQuestionsByQuizID(ctx, session.QuizID)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}

	score := 0
	for _, q := range questions {
		userAnswers := userAnswersByQuestion[q.ID]

		correctAnswers, err := s.answerStore.GetCorrectAnswersByQuestionID(ctx, q.ID) // TODO: měl bych využít MongoDB a i cache
		if err != nil {
			return nil, fmt.Errorf("failed to get correct answers: %w", err)
		}

		correct, err := utils.IsQuestionCorrect(q, userAnswers, correctAnswers)
		if err != nil {
			return nil, err
		}
		if correct {
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

// GetActiveSessionsByQuizID implements service.SessionService.
func (s *sessionService) GetActiveSessionsByQuizID(ctx context.Context, quizID string) ([]*model.Session, error) {
	parsed, err := uuid.Parse(quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session ID: %w", err)
	}
	return s.sessionStore.GetActiveSessionsByQuizID(ctx, parsed)
}
