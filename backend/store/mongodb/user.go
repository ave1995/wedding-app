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
	IconUrl  string `bson:"iconUrl"`
	IsGuest  bool   `bson:"isGuest"`
	QuizID   string `bson:"quizId,omitempty"`
}

const (
	userFieldID      = "_id"
	userFieldEmail   = "email"
	userFieldIsGuest = "isGuest"
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
