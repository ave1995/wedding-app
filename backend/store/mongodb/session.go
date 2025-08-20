package mongodb

import (
	"fmt"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type session struct {
	ID            string `bson:"_id"`
	UserID        string `bson:"user_id"`
	QuizID        string `bson:"quiz_id"`
	CurrentQIndex int64  `bson:"current_qindex"`
	TotalQCount   int64  `bson:"total_qcount"`
	IsCompleted   bool   `bson:"is_completed"`
}

const (
	SessionFieldID            = FieldID
	SessionFieldUserID        = "user_id"
	SessionFieldQuizID        = "quiz_id"
	SessionFieldIsCompleted   = "is_completed"
	SessionFieldCurrentQIndex = "current_qindex"
)

func (m *session) ToDomain() (*model.Session, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse session ID %q: %w", m.ID, err)
	}

	userID, err := uuid.Parse(m.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID %q: %w", m.UserID, err)
	}

	quizID, err := uuid.Parse(m.QuizID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quiz ID %q: %w", m.QuizID, err)
	}

	return &model.Session{
		ID:            id,
		UserID:        userID,
		QuizID:        quizID,
		CurrentQIndex: m.CurrentQIndex,
		TotalQCount:   m.TotalQCount,
		IsCompleted:   m.IsCompleted,
	}, nil
}
