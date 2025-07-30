package mongodb

import (
	"wedding-app/models/domain"

	"github.com/google/uuid"
)

type MongoUser struct {
	ID          uuid.UUID `bson:"_id"`
	Username    string    `bson:"username"`
	Email       string    `bson:"email,omitempty"`
	IsTemporary bool      `bson:"isTemporary"`
}

func (m *MongoUser) ToDomain() *domain.User {
	return &domain.User{
		ID:          m.ID,
		Username:    m.Username,
		Email:       m.Email,
		IsTemporary: m.IsTemporary,
	}
}
