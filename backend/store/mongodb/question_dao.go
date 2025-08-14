package mongodb

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/domain/store"

	"github.com/google/uuid"
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
		ID:     uuid.NewString(),
		QuizID: quizID.String(),
		Text:   text,
	}
	return createAndConvert(ctx, s.questionsCollection(), mongoQuestion)
}

// GetQuestionByID implements store.QuestionStore.
func (s *questionStore) GetQuestionByID(ctx context.Context, id uuid.UUID) (*model.Question, error) {
	panic("unimplemented")
}

// GetQuestionsByQuizID implements store.QuestionStore.
func (s *questionStore) GetQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]*model.Question, error) {
	panic("unimplemented")
}
