package factory

import (
	"context"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/service/question"
	"wedding-app/store/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

func (f *Factory) QuestionStore(ctx context.Context) (store.QuestionStore, error) {
	f.questionStoreOnce.Do(func() {
		var db *mongo.Database
		db, f.questionStoreErr = f.MongoDatabase(ctx)
		if f.questionStoreErr != nil {
			return
		}
		f.questionStore = mongodb.NewQuestionStore(db)
	})
	return f.questionStore, f.questionStoreErr
}

func (f *Factory) QuestionService(ctx context.Context) (service.QuestionService, error) {
	f.questionServiceOnce.Do(func() {
		var store store.QuestionStore
		store, f.questionServiceErr = f.QuestionStore(ctx)
		if f.questionServiceErr != nil {
			return
		}
		f.questionService = question.NewQuestionService(store)
	})
	return f.questionService, f.questionServiceErr
}
