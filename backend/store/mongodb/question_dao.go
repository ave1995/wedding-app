package mongodb

import (
	"context"
	"fmt"
	"time"
	"wedding-app/domain/model"
	"wedding-app/domain/store"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type questionStore struct {
	database *mongo.Database
}

func NewQuestionStore(database *mongo.Database) store.QuestionStore {
	return &questionStore{database: database}
}

const QuestionsCollection = "questions"

func (s *questionStore) questionsCollection() *mongo.Collection {
	return s.database.Collection(QuestionsCollection)
}

// CreateQuestion implements store.QuestionStore.
func (s *questionStore) CreateQuestion(ctx context.Context, text string, quizID uuid.UUID) (*model.Question, error) {
	mongoQuestion := &question{
		ID:        uuid.NewString(),
		QuizID:    quizID.String(),
		Text:      text,
		CreatedAt: time.Now(),
	}
	return createAndConvert(ctx, s.questionsCollection(), mongoQuestion)
}

// GetQuestionByID implements store.QuestionStore.
func (s *questionStore) GetQuestionByID(ctx context.Context, id uuid.UUID) (*model.Question, error) {
	return getByFilterAndConvert[*question](ctx, s.questionsCollection(), bson.M{QuestionFieldID: id.String()})
}

// GetQuestionByIDAndQuizID implements store.QuestionStore.
func (s *questionStore) GetQuestionByIDAndQuizID(ctx context.Context, questionID uuid.UUID, quizID uuid.UUID) (*model.Question, error) {
	return getByFilterAndConvert[*question](ctx, s.questionsCollection(), bson.M{QuestionFieldID: questionID.String(), QuestionFieldQuizID: quizID.String()})
}

// GetQuestionsByQuizID implements store.QuestionStore.
func (s *questionStore) GetQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]*model.Question, error) {
	return getManyByFilterAndConvert[*question](ctx, s.questionsCollection(), bson.M{QuestionFieldQuizID: quizID.String()}, nil)
}

// GetOrderedQuestionsByQuizID implements store.QuestionStore.
func (s *questionStore) GetOrderedQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]*model.Question, error) {
	return getManyByFilterAndConvert[*question](ctx, s.questionsCollection(), bson.M{QuestionFieldQuizID: quizID.String()}, &bson.D{{Key: "created_at", Value: 1}})
}

// GetCountQuestionsByQuizID implements store.QuestionStore.
func (s *questionStore) GetCountQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) (int64, error) {
	filter := bson.M{QuestionFieldQuizID: quizID.String()}

	count, err := s.questionsCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count questions: %w", err)
	}

	return count, nil
}
