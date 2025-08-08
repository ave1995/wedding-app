package factory

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"wedding-app/api/restapi"
	"wedding-app/config"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
	"wedding-app/service/jwt"
	"wedding-app/service/quiz"
	"wedding-app/service/svg"
	"wedding-app/service/user"
	googlecloud "wedding-app/store/googleCloud"
	"wedding-app/store/mongodb"

	"cloud.google.com/go/storage"
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

	googleCloudClient *storage.Client
	svgStore          store.SvgStore
	svgStoreOnce      sync.Once
	svgStoreErr       error

	svgService     service.SvgService
	svgServiceOnce sync.Once
	svgServiceErr  error

	jwtService     service.JWTService
	jwtServiceOnce sync.Once

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
		f.userService = user.NewUserService(store, f.JWTService())
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

func (f *Factory) SvgStore(ctx context.Context) (store.SvgStore, error) {
	f.svgStoreOnce.Do(func() {
		var client *storage.Client
		client, f.svgStoreErr = googlecloud.ConnectClient(ctx)
		if f.svgStoreErr != nil {
			return
		}

		f.svgStore = googlecloud.NewCloud(client, f.config.BucketConfig(), f.Logger())
	})
	return f.svgStore, f.svgStoreErr
}

func (f *Factory) SvgService(ctx context.Context) (service.SvgService, error) {
	f.svgServiceOnce.Do(func() {
		store, err := f.SvgStore(ctx)
		if err != nil {
			f.svgServiceErr = fmt.Errorf("failed to initialize svg service: %w", err)
			return
		}
		f.svgService = svg.NewSvgService(store)
	})
	return f.svgService, f.svgServiceErr
}

func (f *Factory) JWTService() service.JWTService {
	f.jwtServiceOnce.Do(func() {
		f.jwtService = jwt.NewJWTService(f.config.AuthConfig(), f.Logger())
	})
	return f.jwtService
}

func (f *Factory) GinHandlers(ctx context.Context) (*restapi.GinHandlers, error) {
	f.ginHandlersOnce.Do(func() {
		var svgService service.SvgService
		svgService, f.ginHandlersError = f.SvgService(ctx)
		if f.ginHandlersError != nil {
			return
		}
		basicHandler := restapi.NewBasicHandler(svgService)

		var userService service.UserService
		userService, f.ginHandlersError = f.UserService(ctx)
		if f.ginHandlersError != nil {
			return
		}
		userHandler := restapi.NewUserHandler(userService, f.Logger())

		var quizService service.QuizService
		quizService, f.ginHandlersError = f.QuizService(ctx)
		if f.ginHandlersError != nil {
			return
		}
		quizHandler := restapi.NewQuizHandler(quizService)

		authMiddleware := restapi.AuthMiddleware(f.JWTService())

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
