package mongodb

import (
	"context"

	"wedding-app/domain/model"
	"wedding-app/domain/store"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type answerStore struct {
	database *mongo.Database
}

func NewAnswerStore(database *mongo.Database) store.AnswerStore {
	return &answerStore{database: database}
}

// RADEK: malé písmeno, nepotřebuješ exportovat
const AnswerCollection = "answers"

func (s *answerStore) answersCollection() *mongo.Collection {
	return s.database.Collection(AnswerCollection)
}

// CreateAnswer implements store.AnswerStore.
func (s *answerStore) CreateAnswer(ctx context.Context, text string, questionID uuid.UUID, isCorrect bool) (*model.Answer, error) {
	mongoAnswer := &answer{
		ID:         uuid.NewString(),
		QuestionID: questionID.String(),
		Text:       text,
		IsCorrect:  isCorrect,
	}
	// RADEK: mě asi nepříjde nic špatného na tom, kdyby byla funkce create, klidně generická a result by byl answers, nad kterým by si následně zavolal toDomain()
	// Tahle funkce a ani žádná z těch domain_mapperu (getByFilterAndConvert, getByFilterAndConvert, getManyByFilterAndConvert, getManyByFilterAndConvert) nedodržuje single responsibility
	// Když budeš chtít vytvořit odpověd, a ještě s tím něco udělat tak máš smulu protože už si to převedl do modelu
	return createAndConvert(ctx, s.answersCollection(), mongoAnswer)
}

// GetAnswerByID implements store.AnswerStore.
func (s *answerStore) GetAnswerByID(ctx context.Context, id uuid.UUID) (*model.Answer, error) {
	// RADEK: jinak tyhle generické funkce se mi líbí, ale vypadájí generované AI :D
	return getByFilterAndConvert[*answer](ctx, s.answersCollection(), bson.M{AnswerFieldID: id.String()})
}

// GetAnswerByIDAndQuestionID implements store.AnswerStore.
func (s *answerStore) GetAnswerByIDAndQuestionID(ctx context.Context, answerID uuid.UUID, questionID uuid.UUID) (*model.Answer, error) {
	return getByFilterAndConvert[*answer](ctx, s.answersCollection(), bson.M{AnswerFieldID: answerID.String(), AnswerFieldQuestionID: questionID.String()})
}

// GetAnswersByQuestionID implements store.AnswerStore.
func (s *answerStore) GetAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*model.Answer, error) {
	return getManyByFilterAndConvert[*answer](ctx, s.answersCollection(), bson.M{AnswerFieldQuestionID: questionID.String()}, nil)
}

// GetCorrectAnswersByQuestionID implements store.AnswerStore.
func (s *answerStore) GetCorrectAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]*model.Answer, error) {
	return getManyByFilterAndConvert[*answer](ctx, s.answersCollection(), bson.M{AnswerFieldQuestionID: questionID.String(), AnswerFieldIsCorrect: true}, nil)
}
