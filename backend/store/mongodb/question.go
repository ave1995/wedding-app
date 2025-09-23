package mongodb

import (
	"fmt"
	"time"

	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type question struct {
	ID        string    `bson:"_id"`
	QuizID    string    `bson:"quiz_id"`
	Text      string    `bson:"text"`
	Type      string    `bson:"type"`
	PhotoPath *string   `bson:"photo_path,omitempty"`
	CreatedAt time.Time `bson:"created_at"`
}

const (
	QuestionFieldID        = FieldID
	QuestionFieldQuizID    = "quiz_id"
	QuestionFieldCreatedAt = "created_at"
)

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
		ID:        id,
		QuizID:    quizID,
		Text:      q.Text,
		CreatedAt: q.CreatedAt,
		// RADEK: co se stane, když v DB někdo změní typ na něco co neznáš?
		// Tady by si měl v modelu pro QuestionType mít validaci, že ti tam nepříjde něco co neznáš.
		Type:      model.QuestionType(q.Type),
		PhotoPath: q.PhotoPath,
	}, nil
}
