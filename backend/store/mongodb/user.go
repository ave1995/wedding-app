package mongodb

import (
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type user struct {
	ID          string `bson:"_id"`
	Username    string `bson:"username"`
	Email       string `bson:"email,omitempty"`
	IsTemporary bool   `bson:"isTemporary"`
	Password    string `bson:"password"`
}

func (m *user) ToDomain() (*model.User, error) {
	id, err := uuid.Parse(m.ID)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:          id,
		Username:    m.Username,
		Email:       m.Email,
		IsTemporary: m.IsTemporary,
	}, nil
}
