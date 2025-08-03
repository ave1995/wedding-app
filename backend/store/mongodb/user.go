package mongodb

import (
	"wedding-app/domain/model"

	"github.com/google/uuid"
)

type user struct {
	ID          uuid.UUID `bson:"_id"`
	Username    string    `bson:"username"`
	Email       string    `bson:"email,omitempty"`
	IsTemporary bool      `bson:"isTemporary"`
	Password    string    `bson:"password"`
}

func (m *user) ToDomain() *model.User {
	return &model.User{
		ID:          m.ID,
		Username:    m.Username,
		Email:       m.Email,
		IsTemporary: m.IsTemporary,
	}
}
