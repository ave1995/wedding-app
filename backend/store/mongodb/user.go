package mongodb

import (
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type user struct {
	ID       string `bson:"_id"`
	Username string `bson:"username"`
	Email    string `bson:"email,omitempty"`
	Password string `bson:"password"`
	IconUrl  string `bson:"icon_url"`
	IsGuest  bool   `bson:"is_guest"`
	QuizID   string `bson:"quiz_id,omitempty"`
}

const (
	UserFieldID      = FieldID
	UserFieldEmail   = "email"
	UserFieldIsGuest = "is_guest"
)

func (m *user) ToDomain() (*model.User, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, err
	}

	var quizID uuid.UUID
	if m.QuizID != "" {
		parsedQuizID, err := uuid.Parse(m.QuizID)
		if err != nil {
			return nil, err
		}
		quizID = parsedQuizID
	}

	return &model.User{
		ID:       id,
		Username: m.Username,
		Email:    m.Email,
		IsGuest:  m.IsGuest,
		IconUrl:  m.IconUrl,
		QuizID:   quizID,
	}, nil
}
