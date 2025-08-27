package mongodb

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/domain/store"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type attemptStore struct {
	database *mongo.Database
}

func NewAttemptStore(database *mongo.Database) store.AttemptStore {
	return &attemptStore{database: database}
}

const AttemptCollection = "attempts"

func (s *attemptStore) attemptCollection() *mongo.Collection {
	return s.database.Collection(AttemptCollection)
}

// CreateAttemptAnswer implements store.AttemptStore.
func (s *attemptStore) CreateAttemptAnswer(ctx context.Context, params model.CreateAttemptParams) (*model.Attempt, error) {
	mongoAttempt := &attempt{
		ID:         uuid.NewString(),
		SessionID:  params.SessionID.String(),
		QuestionID: params.QuestionID.String(),
		AnswerID:   params.AnswerID.String(),
		IsCorrect:  params.IsCorrect,
	}
	return createAndConvert(ctx, s.attemptCollection(), mongoAttempt)
}

// GetAnsweredBySessionIDAndQuestionID implements store.AttemptStore.
func (s *attemptStore) GetAnsweredBySessionIDAndQuestionID(ctx context.Context, sessionID uuid.UUID, questionID uuid.UUID) (*model.Attempt, error) {
	return getByFilterAndConvert[*attempt](ctx, s.attemptCollection(), bson.M{AttemptFieldSessionID: sessionID.String(), AttemptFieldQuestionID: questionID.String()})

}
