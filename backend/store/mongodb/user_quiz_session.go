package mongodb

import (
	"fmt"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type userQuizSession struct {
	ID          string `bson:"_id"`
	UserID      string `bson:"user_id"`
	QuizID      string `bson:"quiz_id"`
	CurrentQID  string `bson:"current_qid"`
	IsCompleted bool   `bson:"is_completed"`
}

func (m *userQuizSession) ToDomain() (*model.UserQuizSession, error) {
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

	currentQID, err := uuid.Parse(m.CurrentQID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current question ID %q: %w", m.CurrentQID, err)
	}

	return &model.UserQuizSession{
		ID:          id,
		UserID:      userID,
		QuizID:      quizID,
		CurrentQID:  currentQID,
		IsCompleted: m.IsCompleted,
	}, nil
}
