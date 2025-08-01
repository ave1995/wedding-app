package main

import (
	"context"
	"fmt"
	"os"
	"wedding-app/api/restapi"
	"wedding-app/config"
	"wedding-app/factory"
	"wedding-app/utils"

	_ "wedding-app/docs"
)

func main() {
	fmt.Printf("Hello world... %s", os.Getenv("TEST_ENV"))

	config, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	factory := factory.NewFactory(config)

	context := context.Background()

	logger := factory.Logger()

	userService, err := factory.UserService(context)
	if err != nil {
		logger.Error("failed to initialize user service: %w", utils.ErrAttr(err))
	}

	basicHandler := restapi.NewBasicHandler()
	userHandler := restapi.NewUserHandler(userService)

	allHandlers := &restapi.Handlers{
		User:  userHandler,
		Basic: basicHandler,
	}

	server := restapi.NewApiServer(allHandlers, logger, config.ServerConfig())
	logger.Info("starting server", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to listen and server: %w", utils.ErrAttr(err))
	}
}
