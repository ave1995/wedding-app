package server

import (
	"log/slog"
	"wedding-app/handlers"
	"wedding-app/internal/utils"
	"wedding-app/repositories/mongodb"
	"wedding-app/services"
	"wedding-app/storage/mongo"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// swagger embed files

func registerRoutes(r *gin.Engine, config *utils.Configuration) {

	mongoClient, err := mongo.ConnectClient(config.Database.Url)
	if err != nil {
		slog.Error("Error occurred", ErrAttr(err))
	}
	userRepository := mongodb.NewUserRepository(mongo.GetDatabase(mongoClient, config.Database.DbName))
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", userHandler.RegisterUserHandler)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
