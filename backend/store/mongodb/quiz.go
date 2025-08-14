package mongodb

import (
	"fmt"
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type quiz struct {
	ID         string `bson:"_id"`
	Name       string `bson:"name"`
	InviteCode string `bson:"inviteCode"`
}

const (
	QuizFieldID         = FieldID
	QuizFieldName       = "name"
	QuizFieldInviteCode = "inviteCode"
)

func (m *quiz) ToDomain() (*model.Quiz, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quiz ID %q: %w", m.ID, err)
	}

	inviteCode, err := uuid.Parse(m.InviteCode)
	if err != nil {
		return nil, fmt.Errorf("failed to parse invite code %q: %w", m.InviteCode, err)
	}

	return &model.Quiz{
		ID:         id,
		Name:       m.Name,
		InviteCode: inviteCode,
	}, nil
}
