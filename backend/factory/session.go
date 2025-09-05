package factory

import (
	"context"
	"wedding-app/assembler"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/service/session"
	"wedding-app/store/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

func (f *Factory) AttemptStore(ctx context.Context) (store.AttemptStore, error) {
	f.attemptStoreOnce.Do(func() {
		var db *mongo.Database
		db, f.attemptStoreErr = f.MongoDatabase(ctx)
		if f.attemptStoreErr != nil {
			return
		}
		f.attemptStore = mongodb.NewAttemptStore(db)
	})
	return f.attemptStore, f.attemptStoreErr
}

func (f *Factory) SessionStore(ctx context.Context) (store.SessionStore, error) {
	f.sessionStoreOnce.Do(func() {
		var db *mongo.Database
		db, f.sessionStoreErr = f.MongoDatabase(ctx)
		if f.sessionStoreErr != nil {
			return
		}
		f.sessionStore = mongodb.NewSessionStore(db)
	})
	return f.sessionStore, f.sessionStoreErr
}

func (f *Factory) SessionService(ctx context.Context) (service.SessionService, error) {
	f.sessionServiceOnce.Do(func() {
		var sessionStore store.SessionStore
		sessionStore, f.sessionServiceErr = f.SessionStore(ctx)
		if f.sessionServiceErr != nil {
			return
		}
		var questionStore store.QuestionStore
		questionStore, f.sessionServiceErr = f.QuestionStore(ctx)
		if f.sessionServiceErr != nil {
			return
		}
		var attempStore store.AttemptStore
		attempStore, f.sessionServiceErr = f.AttemptStore(ctx)
		if f.sessionServiceErr != nil {
			return
		}
		var answerStore store.AnswerStore
		answerStore, f.sessionServiceErr = f.AnswerStore(ctx)
		if f.sessionServiceErr != nil {
			return
		}
		var userStore store.UserStore
		userStore, f.sessionServiceErr = f.UserStore(ctx)
		if f.sessionServiceErr != nil {
			return
		}
		var assembler *assembler.Assembler
		assembler, f.sessionServiceErr = f.Assembler(ctx)
		if f.sessionServiceErr != nil {
			return
		}

		f.sessionService = session.NewSessionService(sessionStore, questionStore, attempStore, answerStore, userStore, assembler, f.EventPublisher(), f.Logger())
	})
	return f.sessionService, f.sessionServiceErr
}
