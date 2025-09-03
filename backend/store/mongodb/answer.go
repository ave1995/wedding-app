package mongodb

import (
	"fmt"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type answer struct {
	ID         string `bson:"_id"`
	QuestionID string `bson:"question_id"`
	Text       string `bson:"text"`
	IsCorrect  bool   `bson:"is_correct"`
}

const (
	AnswerFieldID         = FieldID
	AnswerFieldQuestionID = "question_id"
	AnswerFieldIsCorrect  = "is_correct"
)

func (a *answer) ToDomain() (*model.Answer, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse answer ID %q: %w", a.ID, err)
	}

	questionID, err := uuid.Parse(a.QuestionID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse question ID %q: %w", a.QuestionID, err)
	}

	return &model.Answer{
		ID:         id,
		QuestionID: questionID,
		Text:       a.Text,
		IsCorrect:  a.IsCorrect,
	}, nil
}
