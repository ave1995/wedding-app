package mongodb

import (
	"context"
	"fmt"
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
	collection := s.quizzesCollection()

	mongoQuiz := &quiz{
		ID:         uuid.NewString(),
		Name:       name,
		InviteCode: uuid.NewString(),
	}

	_, err := collection.InsertOne(ctx, mongoQuiz)
	if err != nil {
		return nil, fmt.Errorf("failed to insert one to: %w", err)
	}

	quiz, err := mongoQuiz.ToDomain()
	if err != nil {
		return nil, err
	}

	return quiz, nil
}

// GetQuizByInviteCode implements store.QuizStore.
func (s *quizStore) GetQuizByInviteCode(ctx context.Context, inviteCode uuid.UUID) (*model.Quiz, error) {
	return s.getQuizByFilter(ctx, bson.M{quizFieldInviteCode: inviteCode.String()})
}

// GetQuizByID implements store.QuizStore.
func (s *quizStore) GetQuizByID(ctx context.Context, id uuid.UUID) (*model.Quiz, error) {
	return s.getQuizByFilter(ctx, bson.M{quizFieldID: id.String()})
}

// getQuizByFilter is a helper that encapsulates shared quiz get logic.
func (s *quizStore) getQuizByFilter(ctx context.Context, filter bson.M) (*model.Quiz, error) {
	result, err := getByFilter[quiz](ctx, s.quizzesCollection(), filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find quiz: %w", err)
	}

	quiz, err := result.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert quiz to domain model: %w", err)
	}

	return quiz, nil
}
