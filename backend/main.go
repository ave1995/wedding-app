package main

import (
	"context"
	"log/slog"
	"wedding-app/api/restapi"
	"wedding-app/config"
	"wedding-app/factory"
	"wedding-app/utils"

	_ "wedding-app/api/restapi/docs"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	factory := factory.NewFactory(config)

	mainContext := context.Background()

	logger := factory.Logger()

	userService, err := factory.UserService(mainContext)
	if err != nil {
		logger.Error("failed to initialize user service: ", utils.ErrAttr(err))
	}

	basicHandler := restapi.NewBasicHandler()
	userHandler := restapi.NewUserHandler(userService)

	allHandlers := &restapi.Handlers{
		User:  userHandler,
		Basic: basicHandler,
	}

	server := restapi.NewApiServer(allHandlers, logger, config.ServerConfig())
	logger.InfoContext(context.Background(), "starting server", slog.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to listen and server: %w", utils.ErrAttr(err))
	}
}
