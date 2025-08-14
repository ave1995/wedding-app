package mongodb

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/domain/store"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type quizStore struct {
	database *mongo.Database
}

func NewQuizStore(database *mongo.Database) store.QuizStore {
	return &quizStore{database: database}
}

const QuizzesCollection = "quizzes"

func (s *quizStore) quizzesCollection() *mongo.Collection {
	return s.database.Collection(QuizzesCollection)
}

// CreateQuiz implements store.QuizStore.
func (s *quizStore) CreateQuiz(ctx context.Context, name string) (*model.Quiz, error) {
	mongoQuiz := &quiz{
		ID:         uuid.NewString(),
		Name:       name,
		InviteCode: uuid.NewString(),
	}
	return createAndConvert(ctx, s.quizzesCollection(), mongoQuiz)
}

// GetQuizByInviteCode implements store.QuizStore.
func (s *quizStore) GetQuizByInviteCode(ctx context.Context, inviteCode uuid.UUID) (*model.Quiz, error) {
	return getByFilterAndConvert[*quiz](ctx, s.quizzesCollection(), bson.M{QuizFieldInviteCode: inviteCode.String()})
}

// GetQuizByID implements store.QuizStore.
func (s *quizStore) GetQuizByID(ctx context.Context, id uuid.UUID) (*model.Quiz, error) {
	return getByFilterAndConvert[*quiz](ctx, s.quizzesCollection(), bson.M{QuizFieldID: id.String()})
}
