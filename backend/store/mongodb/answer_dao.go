package mongodb

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/domain/store"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type answerStore struct {
	database *mongo.Database
}

func NewAnswerStore(database *mongo.Database) store.AnswerStore {
	return &answerStore{database: database}
}

const AnswerCollection = "answers"

func (s *answerStore) answersCollection() *mongo.Collection {
	return s.database.Collection(AnswerCollection)
}

// CreateAnswer implements store.AnswerStore.
func (s *answerStore) CreateAnswer(ctx context.Context, text string, questionID uuid.UUID) (*model.Answer, error) {
	panic("unimplemented")
}

// GetAnswerByID implements store.AnswerStore.
func (s *answerStore) GetAnswerByID(ctx context.Context, id uuid.UUID) (*model.Answer, error) {
	panic("unimplemented")
}

// GetAnswersByQuestionID implements store.AnswerStore.
func (s *answerStore) GetAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*model.Answer, error) {
	panic("unimplemented")
}
