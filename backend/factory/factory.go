package factory

import (
	"context"
	"fmt"
	"log/slog"
	"wedding-app/config"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/store/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

type Factory struct {
	config config.Config

	logger *slog.Logger

	mongoClient   *mongo.Client
	mongoDatabase *mongo.Database

	userStore   store.UserStore
	userService service.UserService
}

func NewFactory(config config.Config) *Factory {
	return &Factory{
		config: config,
	}
}

func (f *Factory) Logger() *slog.Logger {
	if f.logger == nil {
		f.logger = getLogger()
	}
	return f.logger
}

func (f *Factory) MongoClient(ctx context.Context) (*mongo.Client, error) {
	if f.mongoClient != nil {
		return f.mongoClient, nil
	}

	var err error

	f.mongoClient, err = getMongoClient(ctx, f.Logger(), f.config.StoreConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	return f.mongoClient, nil
}

func (f *Factory) MongoDatabase(ctx context.Context) (*mongo.Database, error) {
	if f.mongoDatabase != nil {
		return f.mongoDatabase, nil
	}

	client, err := f.MongoClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get mongo client: %w", err)
	}

	f.mongoDatabase = mongodb.GetDatabase(client, f.config.StoreConfig().DbName)
	return f.mongoDatabase, err
}

func (f *Factory) UserStore(ctx context.Context) (store.UserStore, error) {
	if f.userStore != nil {
		return f.userStore, nil
	}

	var err error

	database, err := f.MongoDatabase(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get mongo database: %w", err)
	}

	f.userStore = getUserStore(database)
	return f.userStore, nil
}

func (f *Factory) UserService(ctx context.Context) (service.UserService, error) {
	if f.userService != nil {
		return f.userService, nil
	}

	var err error

	store, err := f.UserStore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user service: %w", err)
	}

	f.userService = getUserService(store)
	return f.userService, nil
}
