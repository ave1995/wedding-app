package restapi

import (
	"log/slog"
	"net/http"
	"time"
	"wedding-app/config"

	"github.com/gin-gonic/gin"
)

func NewGinServer(handlers *GinHandlers, logger *slog.Logger, config config.ServerConfig) *http.Server {
	router := gin.New()

	router.Use(gin.Logger())

	//Use my error middleware
	router.Use(ErrorHandler())
	router.Use(gin.Recovery())

	handlers.RegisterAll(router)

	s := &http.Server{
		Addr:           config.Host + ":" + config.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}
