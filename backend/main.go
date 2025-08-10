package main

import (
	"context"
	"log/slog"
	"os"

	"wedding-app/config"
	"wedding-app/factory"
	"wedding-app/utils"

	_ "wedding-app/api/restapi/docs"
)

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_token
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	fac := factory.NewFactory(cfg)

	mainContext := context.Background()

	logger := fac.Logger()

	server, err := fac.HttpServer(mainContext)
	if err != nil {
		logger.Error("failed to initialize http server: ", utils.ErrAttr(err))
		os.Exit(1)
	}

	logger.Info("starting server", slog.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to listen and server: %w", utils.ErrAttr(err))
		fac.Cleanup()
		os.Exit(1)
	}
}
