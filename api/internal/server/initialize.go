package server

import (
	"fmt"
	"log/slog"
	"wedding-app/internal/utils"

	"github.com/gin-gonic/gin"
)

func InitServer(config *utils.Configuration) {

	slog.Info("Configuration information",
		slog.String("mongo_url", config.Database.Url),
		slog.String("server_port", config.Server.Port),
		slog.String("db_name", config.Database.DbName),
	)

	slog.Info(fmt.Sprintf("Application Name %s is starting....", config.App.Name))

	router := gin.Default()
	registerRoutes(router, config)

	formattedUrl := fmt.Sprintf(": %s", config.Server.Port)
	router.Run(formattedUrl)
}
