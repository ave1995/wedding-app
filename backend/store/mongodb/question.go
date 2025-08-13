package mongodb

import (
	"fmt"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type question struct {
	ID     string `bson:"_id"`
	QuizID string `bson:"quizId"`
	Text   string `bson:"text"`
}

func (q *question) ToDomain() (*model.Question, error) {
	id, err := uuid.Parse(q.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse question ID %q: %w", q.ID, err)
	}

	quizID, err := uuid.Parse(q.QuizID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quiz ID %q: %w", q.QuizID, err)
	}

	return &model.Question{
		ID:     id,
		QuizID: quizID,
		Text:   q.Text,
	}, nil
}
