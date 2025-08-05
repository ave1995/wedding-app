package factory

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"wedding-app/api/restapi"
	"wedding-app/config"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/service/user"
	"wedding-app/store/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
)

type Factory struct {
	config config.Config

	logger     *slog.Logger
	loggerOnce sync.Once

	mongoClient   *mongo.Client
	mongoDatabase *mongo.Database
	mongoOnce     sync.Once
	mongoError    error

	userStore     store.UserStore
	userStoreOnce sync.Once
	userStoreErr  error

	userService     service.UserService
	userServiceOnce sync.Once
	userServiceErr  error

	ginHandlers      *restapi.GinHandlers
	ginHandlersOnce  sync.Once
	ginHandlersError error

	server      *http.Server
	serverOnce  sync.Once
	serverError error
}

func NewFactory(config config.Config) *Factory {
	return &Factory{
		config: config,
	}
}

func (f *Factory) Logger() *slog.Logger {
	f.loggerOnce.Do(func() {
		f.logger = newLogger()
	})
	return f.logger
}

func (f *Factory) MongoDatabase(ctx context.Context) (*mongo.Database, error) {
	f.mongoOnce.Do(func() {
		f.mongoClient, f.mongoError = mongodb.ConnectClient(ctx, f.Logger(), f.config.StoreConfig())
		if f.mongoError != nil {
			return
		}
		f.mongoDatabase = mongodb.GetDatabase(f.mongoClient, f.config.StoreConfig().DbName)
	})

	return f.mongoDatabase, f.mongoError
}

func (f *Factory) UserStore(ctx context.Context) (store.UserStore, error) {
	f.userStoreOnce.Do(func() {
		var db *mongo.Database
		db, f.userStoreErr = f.MongoDatabase(ctx)
		if f.userStoreErr != nil {
			return
		}
		f.userStore = mongodb.NewUserStore(db, f.Logger())
	})
	return f.userStore, f.userStoreErr
}

func (f *Factory) UserService(ctx context.Context) (service.UserService, error) {
	f.userServiceOnce.Do(func() {
		var store store.UserStore
		store, f.userServiceErr = f.UserStore(ctx)
		if f.userServiceErr != nil {
			return
		}
		f.userService = user.NewUserService(store)
	})
	return f.userService, f.userServiceErr
}

func (f *Factory) GinHandlers(ctx context.Context) (*restapi.GinHandlers, error) {
	f.ginHandlersOnce.Do(func() {
		basicHandler := restapi.NewBasicHandler()

		var userService service.UserService
		userService, f.ginHandlersError = f.UserService(ctx)
		userHandler := restapi.NewUserHandler(userService)

		f.ginHandlers = &restapi.GinHandlers{
			User:  userHandler,
			Basic: basicHandler,
		}
	})
	return f.ginHandlers, f.ginHandlersError
}

func (f *Factory) HttpServer(ctx context.Context) (*http.Server, error) {
	f.serverOnce.Do(func() {
		var ginHandlers *restapi.GinHandlers
		ginHandlers, f.serverError = f.GinHandlers(ctx)

		f.server = restapi.NewGinServer(ginHandlers, f.Logger(), f.config.ServerConfig())
	})
	return f.server, f.serverError
}
