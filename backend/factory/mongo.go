package factory

import (
	"context"
	"log/slog"
	"sync"
	"wedding-app/config"
	"wedding-app/store/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mongoOnce   sync.Once
	mongoClient *mongo.Client
)

func getMongoClient(ctx context.Context, logger *slog.Logger, config config.StoreConfig) (*mongo.Client, error) {
	var err error
	mongoOnce.Do(func() {
		mongoClient, err = mongodb.ConnectClient(ctx, logger, config.DbUrl)
	})

	return mongoClient, err
}
