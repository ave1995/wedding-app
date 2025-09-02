package assembler

import (
	"context"
	"time"
	"wedding-app/domain/event"
	"wedding-app/domain/store"

	"github.com/google/uuid"
)

type Assembler struct {
	quizStore     store.QuizStore
	userStore     store.UserStore
	questionStore store.QuestionStore
	answerStore   store.AnswerStore
}

func NewAssembler(qs store.QuizStore, us store.UserStore, ques store.QuestionStore, as store.AnswerStore) *Assembler {
	return &Assembler{quizStore: qs, userStore: us, questionStore: ques, answerStore: as}
}

func (a *Assembler) buildSessionBase(ctx context.Context, sessionID, userID, quizID uuid.UUID) (event.SessionBaseEvent, error) {
	quiz, err := a.quizStore.GetQuizByID(ctx, quizID)
	if err != nil {
		return event.SessionBaseEvent{}, err
	}

	user, err := a.userStore.GetUserByID(ctx, userID)
	if err != nil {
		return event.SessionBaseEvent{}, err
	}

	return event.SessionBaseEvent{
		SessionID:   sessionID,
		UserID:      userID,
		Username:    user.Username,
		UserIconUrl: user.IconUrl,
		Quizname:    quiz.Name,
	}, nil
}

func (a *Assembler) ToSessionStartEvent(ctx context.Context, sessionID, userID, quizID uuid.UUID) (*event.SessionStartEvent, error) {
	base, err := a.buildSessionBase(ctx, sessionID, userID, quizID)
	if err != nil {
		return nil, err
	}

	return &event.SessionStartEvent{
		SessionBaseEvent: base,
		StartedAt:        time.Now(),
	}, nil
}

func (a *Assembler) ToSessionEndEvent(ctx context.Context, sessionID, userID, quizID uuid.UUID) (*event.SessionEndEvent, error) {
	base, err := a.buildSessionBase(ctx, sessionID, userID, quizID)
	if err != nil {
		return nil, err
	}

	return &event.SessionEndEvent{
		SessionBaseEvent: base,
		EndedAt:          time.Now(),
	}, nil
}

func (a *Assembler) ToAnswerSubmittedEvent(ctx context.Context, sessionID, userID, quizID, quiestionID uuid.UUID, answerIDs []uuid.UUID) (*event.AnswerSubmittedEvent, error) {
	base, err := a.buildSessionBase(ctx, sessionID, userID, quizID)
	if err != nil {
		return nil, err
	}

	question, err := a.questionStore.GetQuestionByID(ctx, quiestionID)
	if err != nil {
		return nil, err
	}

	answers := make([]string, 0, len(answerIDs))
	for _, id := range answerIDs {
		answer, err := a.answerStore.GetAnswerByID(ctx, id)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer.Text)
	}

	return &event.AnswerSubmittedEvent{
		SessionBaseEvent: base,
		QuestionText:     question.Text,
		Answers:          answers,
		QuestionID:       quiestionID,
		SubmittedAt:      time.Now(),
	}, err
}

func (a *Assembler) ToQuestionOpenedEvent(ctx context.Context, sessionID, userID, quizID, quiestionID uuid.UUID) (*event.QuestionOpenedEvent, error) {
	base, err := a.buildSessionBase(ctx, sessionID, userID, quizID)
	if err != nil {
		return nil, err
	}

	question, err := a.questionStore.GetQuestionByID(ctx, quiestionID)
	if err != nil {
		return nil, err
	}

	return &event.QuestionOpenedEvent{
		SessionBaseEvent: base,
		QuestionText:     question.Text,
		QuestionID:       quiestionID,
		OpenedAt:         time.Now(),
	}, err
}
