package factory

import (
	"context"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/service/answer"
	"wedding-app/store/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

func (f *Factory) AnswerStore(ctx context.Context) (store.AnswerStore, error) {
	f.answerStoreOnce.Do(func() {
		var db *mongo.Database
		db, f.answerStoreErr = f.MongoDatabase(ctx)
		if f.answerStoreErr != nil {
			return
		}
		f.answerStore = mongodb.NewAnswerStore(db)
	})
	return f.answerStore, f.answerStoreErr
}

func (f *Factory) AnswerService(ctx context.Context) (service.AnswerService, error) {
	f.answerServiceOnce.Do(func() {
		var store store.AnswerStore
		store, f.answerServiceErr = f.AnswerStore(ctx)
		if f.answerServiceErr != nil {
			return
		}
		f.anwserService = answer.NewAnswerService(store)
	})
	return f.anwserService, f.answerServiceErr
}
