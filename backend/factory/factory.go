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
	"wedding-app/service/jwt"
	"wedding-app/service/quiz"
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

	quizStore     store.QuizStore
	quizStoreOnce sync.Once
	quizStoreErr  error

	quizService     service.QuizService
	quizServiceOnce sync.Once
	quizServiceErr  error

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

func (f *Factory) QuizStore(ctx context.Context) (store.QuizStore, error) {
	f.quizStoreOnce.Do(func() {
		var db *mongo.Database
		db, f.quizStoreErr = f.MongoDatabase(ctx)
		if f.quizStoreErr != nil {
			return
		}
		f.quizStore = mongodb.NewQuizStore(db)
	})
	return f.quizStore, f.quizStoreErr
}

func (f *Factory) QuizService(ctx context.Context) (service.QuizService, error) {
	f.quizServiceOnce.Do(func() {
		var store store.QuizStore
		store, f.quizServiceErr = f.QuizStore(ctx)
		if f.quizServiceErr != nil {
			return
		}
		f.quizService = quiz.NewQuizService(store)
	})
	return f.quizService, f.quizServiceErr
}

func (f *Factory) GinHandlers(ctx context.Context) (*restapi.GinHandlers, error) {
	f.ginHandlersOnce.Do(func() {
		basicHandler := restapi.NewBasicHandler()

		var userService service.UserService
		userService, f.ginHandlersError = f.UserService(ctx)
		if f.ginHandlersError != nil {
			return
		}
		userHandler := restapi.NewUserHandler(userService)

		var quizService service.QuizService
		quizService, f.ginHandlersError = f.QuizService(ctx)
		if f.ginHandlersError != nil {
			return
		}
		quizHandler := restapi.NewQuizHandler(quizService)

		var jwtService service.JWTService
		jwtService, f.ginHandlersError = jwt.NewJWTService(f.config.AuthConfig(), f.Logger())
		if f.ginHandlersError != nil {
			return
		}

		authMiddleware := restapi.AuthMiddleware(jwtService)

		f.ginHandlers, f.ginHandlersError = restapi.NewGinHandlers(userHandler, basicHandler, quizHandler, authMiddleware)
	})
	return f.ginHandlers, f.ginHandlersError
}

func (f *Factory) HttpServer(ctx context.Context) (*http.Server, error) {
	f.serverOnce.Do(func() {
		var ginHandlers *restapi.GinHandlers
		ginHandlers, f.serverError = f.GinHandlers(ctx)
		if f.serverError != nil {
			return
		}

		f.server = restapi.NewGinServer(ginHandlers, f.Logger(), f.config.ServerConfig())
	})
	return f.server, f.serverError
}
