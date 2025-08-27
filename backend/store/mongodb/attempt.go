package mongodb

import (
	"fmt"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type attempt struct {
	ID         string `bson:"_id"`
	SessionID  string `bson:"session_id"`
	QuestionID string `bson:"question_id"`
	AnswerID   string `bson:"answer_id"`
	IsCorrect  bool   `bson:"is_correct"`
}

const (
	AttemptFieldID         = FieldID
	AttemptFieldSessionID  = "session_id"
	AttemptFieldQuestionID = "question_id"
	AttemptFieldAnswerID   = "answer_id"
)

func (m *attempt) ToDomain() (*model.Attempt, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse attempt answer ID %q: %w", m.ID, err)
	}

	sessionID, err := uuid.Parse(m.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session ID %q: %w", m.SessionID, err)
	}

	questionID, err := uuid.Parse(m.QuestionID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse question ID %q: %w", m.QuestionID, err)
	}

	answerID, err := uuid.Parse(m.AnswerID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse answer ID %q: %w", m.AnswerID, err)
	}

	return &model.Attempt{
		ID:         id,
		SessionID:  sessionID,
		QuestionID: questionID,
		AnswerID:   answerID,
		IsCorrect:  m.IsCorrect,
	}, nil
}
