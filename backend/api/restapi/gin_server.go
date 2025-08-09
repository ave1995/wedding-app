package restapi

import (
	"log/slog"
	"net/http"
	"time"
	"wedding-app/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func NewGinServer(handlers *GinHandlers, logger *slog.Logger, config config.ServerConfig) *http.Server {
	router := gin.New()

	// Add the sloggin middleware to all routes.
	// The middleware will log all requests attributes.
	router.Use(sloggin.New(logger))
	router.Use(gin.Recovery())
	router.Use(ErrorHandler(logger))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     config.Origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
