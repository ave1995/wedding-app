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
func (s *questionStore) CreateQuestion(ctx context.Context, text string, quizID uuid.UUID, questionType model.QuestionType) (*model.Question, error) {
	mongoQuestion := &question{
		ID:        uuid.NewString(),
		QuizID:    quizID.String(),
		Text:      text,
		Type:      string(questionType),
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
func (s *questionStore) GetCountQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) (int, error) {
	filter := bson.M{QuestionFieldQuizID: quizID.String()}

	count64, err := s.questionsCollection().CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count questions: %w", err)
	}

	// Check for overflow on 32-bit systems
	if count64 > int64(^uint(0)>>1) { // max int value
		return 0, fmt.Errorf("question count too large to fit in int: %d", count64)
	}

	return int(count64), nil
}
