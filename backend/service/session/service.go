package session

import (
	"context"
	"fmt"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"

	"github.com/google/uuid"
)

type sessionService struct {
	sessionStore       store.SessionStore
	questionStore      store.QuestionStore
	attemptAnswerStore store.AttemptAnswerStore
	answerStore        store.AnswerStore
}

func NewSessionService(ss store.SessionStore, qs store.QuestionStore, aas store.AttemptAnswerStore, as store.AnswerStore) service.SessionService {
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
	questions, err := s.questionStore.GetQuestionsByQuizID(ctx, session.QuizID)
	if err != nil {
		return nil, fmt.Errorf("failed to load questions: %w", err)
	}
	// Get All Answered Questions
	answered, err := s.attemptAnswerStore.GetAnsweredBySession(ctx, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load attempt answers: %w", err)
	}

	nextQuestion := getNextQuestion(answered, questions)

	// If question is null that means all questions answered and session is completed
	if nextQuestion == nil {
		session.IsCompleted = true
		if err := s.sessionStore.UpdateSession(ctx, session); err != nil {
			return nil, fmt.Errorf("failed to mark session completed: %w", err)
		}
		return nil, nil
	}
	// I update session with current question ID for remembering which question show
	session.CurrentQID = nextQuestion.ID
	if err := s.sessionStore.UpdateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to update current question: %w", err)
	}

	return nextQuestion, nil
}

// StartQuiz implements service.SessionService.
func (s *sessionService) StartQuiz(ctx context.Context, userID string, quizID string) (*model.UserQuizSession, error) {
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
	return s.sessionStore.CreateSession(ctx, parsedUserID, parsedQuizID)
}

// SubmitAnswer implements service.SessionService.
func (s *sessionService) SubmitAnswer(ctx context.Context, sessionID string, questionID string, answerID string) (isCompleted bool, err error) {
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
	session, err := s.sessionStore.FindByID(ctx, parsedSessionID)
	if err != nil {
		return false, fmt.Errorf("failed to load session: %w", err)
	}
	if session.IsCompleted {
		return false, fmt.Errorf("session already completed")
	}
	questions, err := s.questionStore.GetQuestionsByQuizID(ctx, session.QuizID)
	if err != nil {
		return false, fmt.Errorf("failed to load questions: %w", err)
	}
	var question *model.Question
	for _, q := range questions {
		if q.ID == parsedQuestionID {
			question = q
			break
		}
	}
	if question == nil {
		return false, fmt.Errorf("question does not belong to this quiz")
	}
	answers, err := s.answerStore.GetAnswersByQuestionID(ctx, question.ID)
	if err != nil {
		return false, fmt.Errorf("failed to load answers: %w", err)
	}
	isCorrect := false
	for _, ans := range answers {
		if ans.ID == parsedAnswerID {
			isCorrect = ans.IsCorrect
			break
		}
	}

	if _, err := s.attemptAnswerStore.CreateAttemptAnswer(ctx, model.CreateAttemptAnswerParams{
		SessionID:  session.ID,
		QuestionID: question.ID,
		AnswerID:   parsedAnswerID,
		IsCorrect:  isCorrect,
	}); err != nil {
		return false, fmt.Errorf("failed to save answer: %w", err)
	}

	// 7. Find next unanswered question
	answered, err := s.attemptAnswerStore.GetAnsweredBySession(ctx, session.ID)
	if err != nil {
		return false, fmt.Errorf("failed to load attempt answers: %w", err)
	}
	nextQuestion := getNextQuestion(answered, questions)

	// 8. Update session
	if nextQuestion.ID == uuid.Nil {
		session.IsCompleted = true
		session.CurrentQID = uuid.Nil
		isCompleted = true
	} else {
		session.CurrentQID = nextQuestion.ID
	}
	if err := s.sessionStore.UpdateSession(ctx, session); err != nil {
		return false, fmt.Errorf("failed to update session: %w", err)
	}

	return isCompleted, nil
}

func getNextQuestion(answered []*model.AttemptAnswer, questions []*model.Question) *model.Question {
	// I map answered by question ID
	answeredMap := make(map[uuid.UUID]bool)
	for _, a := range answered {
		answeredMap[a.QuestionID] = true
	}
	// I get next question If map does not contain it
	var nextQuestion *model.Question
	for _, q := range questions {
		if !answeredMap[q.ID] {
			nextQuestion = q
			break
		}
	}

	return nextQuestion
}
