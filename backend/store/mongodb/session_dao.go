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

type sessionStore struct {
	database *mongo.Database
}

func NewSessionStore(database *mongo.Database) store.SessionStore {
	return &sessionStore{database: database}
}

const SessionCollection = "sessions"

func (s *sessionStore) sessionCollection() *mongo.Collection {
	return s.database.Collection(SessionCollection)
}

// CreateSession implements store.SessionStore.
func (s *sessionStore) CreateSession(ctx context.Context, userID uuid.UUID, quizID uuid.UUID, questionCount int) (*model.Session, error) {
	mongoSession := &session{
		ID:            uuid.NewString(),
		UserID:        userID.String(),
		QuizID:        quizID.String(),
		TotalQCount:   questionCount,
		CurrentQIndex: 0,
		IsCompleted:   false,
	}
	return createAndConvert(ctx, s.sessionCollection(), mongoSession)
}

// FindActive implements store.SessionStore.
func (s *sessionStore) FindActive(ctx context.Context, userID uuid.UUID, quizID uuid.UUID) (*model.Session, error) {
	return getByFilterAndConvert[*session](ctx, s.sessionCollection(), bson.M{SessionFieldUserID: userID.String(), SessionFieldQuizID: quizID.String()})
}

// FindByID implements store.SessionStore.
func (s *sessionStore) FindByID(ctx context.Context, sessionID uuid.UUID) (*model.Session, error) {
	return getByFilterAndConvert[*session](ctx, s.sessionCollection(), bson.M{SessionFieldID: sessionID.String()})
}

// UpdateSession implements store.SessionStore.
func (s *sessionStore) UpdateSession(ctx context.Context, session *model.Session) error {
	filter := bson.M{SessionFieldID: session.ID.String()}

	update := bson.M{
		"$set": bson.M{
			SessionFieldIsCompleted:   session.IsCompleted,
			SessionFieldCurrentQIndex: session.CurrentQIndex,
		},
	}
	fmt.Printf("Type of session.ID: %T, value: %v\n", session.ID, session.ID)

	res, err := s.sessionCollection().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("no session found with id %v", session.ID)
	}

	return nil
}
