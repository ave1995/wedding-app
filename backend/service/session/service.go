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
	userStore          store.UserStore
	assembler          *assembler.Assembler
	publisher          event.EventPublisher
	logger             *slog.Logger
}

func NewSessionService(ss store.SessionStore, qs store.QuestionStore, aas store.AttemptStore, as store.AnswerStore, us store.UserStore, a *assembler.Assembler, pub event.EventPublisher, logger *slog.Logger) service.SessionService {
	return &sessionService{sessionStore: ss, questionStore: qs, attemptAnswerStore: aas, answerStore: as, userStore: us, assembler: a, publisher: pub, logger: logger}
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
	// RADEK: data stojí většinou peníze a čas,
	// Tady bych udělal metodu, která ti vráti konkrétní otázku pro index, pokud neexistuje nebo je mimo range, tak to bych kontroloval před nebo při načítání
	// Samozřejmě potom ti bude chybět počet otázek, což je informace, kterou můžeš vědět při spuštění session a mit jí uloženou tam, nemusíš pokaždé načítat stejné číslo
	// Co když se v průběhu změní počet otázek? To bych rozdoně nepovoloval, to znamená quiz by měl mít nějaké stavy, jakmile se stav změní na "publikovaný" nelze upravovat, pokud není "publikovaný" nelze vyplnovat.
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
	// RADEK: toto je hack aby uživatel neviděl odpovědi co?
	// Takže mít answer model a potom nějaký model, který to oseká pro uživatele a nad answer mít metodu, která udělá convert do nového modelu. žádný DTO :(
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
	// RADEK: dva ify za sebou většinou znaměnají "přepiš to do switche" je to více čitelné potom :D
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
	// RADEK: ta session dokonce už obsahuje kolik je tam otázek :D takže to co jsem psal výše dáva ještě větší smysl.
	newSession, err := s.sessionStore.CreateSession(ctx, parsedUserID, parsedQuizID, count)
	if err != nil {
		return nil, err
	}

	// RADEK: nelíbí se mi tenhle sestavovač :D
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
			IsCorrect:  answer.IsCorrect, // RADEK: Za mě je tohle redundantní informace, zase na druhou stranu tohle se v documents dbs dělá asi no
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

		// RADEK: tady jak mluvíš o te cache, tak to by na urovni service vlastně nemělo být vidět, stačí ti vyměnit kdyžtak ten answerStore
		// Jediné, kde by to mohlo být třeba vidět je při invalidaci cache, takže tam kde by si něco měnil.
		correctAnswers, err := s.answerStore.GetCorrectAnswersByQuestionID(ctx, q.ID) // TODO: měl bych využít MongoDB a i cache
		if err != nil {
			// RADEK: je také takové pravidlo, že error neobsahuje takové ty slova "failed", "error", ...
			// Takže takový error by vypadal takto: "get correct answers"
			return nil, fmt.Errorf("failed to get correct answers: %w", err)
		}

		// RADEK: zbytečný kod, error nikdy nevznikne. potom můžeš krásně přepsat:
		//if correct := utils.IsQuestionCorrect(q, userAnswers, correctAnswers); correct {
		//	score++
		//}

		correct, err := utils.IsQuestionCorrect(q, userAnswers, correctAnswers)
		if err != nil {
			return nil, err
		}
		if correct {
			score++
		}
	}

	// RADEK: tohle bych klidně taky zabalil do utils aby si mohl jednoduše pokrýt testy.
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

// GetResultsByQuizID implements service.SessionService.
func (s *sessionService) GetResultsByQuizID(ctx context.Context, quizID string) ([]*dto.UserResult, error) {
	parsed, err := uuid.Parse(quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session ID: %w", err)
	}

	completed, err := s.sessionStore.GetCompletedSessionsByQuizID(ctx, parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}

	// RADEK: Tady by si asi měl vědět kolik jích bude žejo?
	// userResults := make([]*dto.UserResult, 0, len(completed))
	var userResults []*dto.UserResult
	for _, session := range completed {
		// RADEK: Tohle je zase velká vykřičník = špatný návrh
		// 1. Volat public metodu služby z jiné public metody znamená, že je něco špatně
		// 2. GetResult načítá session podle sessionID, ale session tady už máš takže tady 4 řádky nazpátek načteš 100 sessions z databáze aby si pro každou session zavolal načtení sesssion? (chudák databáze)
		// 3. Načítání probíhá pro quiz ID takže si dopředu můžeš načíst quiz, jeho otázky a odpovědi to je to co pro každé volání dělá GetResult - pokud správně počítám tak tahle funkce
		// by mohla volat tak 5x databázi a potom všechno dopočítat :D za to volá 5 x počet session při 1_000 sessions to je 5_000 :) 5 vs 5_000 - v distribuovaných systémech a v cloudu je nejdražší a nejvíc spomalující věc databáze tak na to mysli
		result, err := s.GetResult(ctx, session.ID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to get result: %w", err)
		}

		// RADEK: tohle se vlastně používalo na vypsání toho žebříčku ne? Tak to se už nedivím, že to trvalo tak dlouho :D
		// Velký vykřičník je vždycky když se store volá v cyklu - to většinou znamená že to pujde udělat líp, v tomhle případě si dopředu jedním dotazem můžeš zavolat všechny uživatele pro ten quiiz nebo pro ty sessionsIDs když si je dáš do pole
		// Takže ty uživatele načteš jedním příkazem, uděláš si z toho mapu do které už tady jen lezeš - úspora obrovská
		user, err := s.userStore.GetUserByID(ctx, session.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}

		// RADEK: Na co je tahle proměnná dobrá?
		// userResults = append(userResults, &dto.UserResult{
		//			Result: result,
		//			User:   user,
		//		})
		userResult := &dto.UserResult{
			Result: result,
			User:   user,
		}

		userResults = append(userResults, userResult)
	}

	// RADEK: Z tehle funkce mám pocit, že si jí dodělával narychlo, den před svatbou ne? :D
	return userResults, nil
}
