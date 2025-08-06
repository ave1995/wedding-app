package restapi

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinHandlers struct {
	User           *UserHandler
	Basic          *BasicHandler
	AuthMiddleware gin.HandlerFunc
}

func NewGinHandlers(user *UserHandler, basic *BasicHandler, authMiddleware gin.HandlerFunc) *GinHandlers {
	return &GinHandlers{
		User:           user,
		Basic:          basic,
		AuthMiddleware: authMiddleware,
	}
}

func (h *GinHandlers) RegisterAll(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	h.User.Register(auth)

	api := router.Group("/api")
	api.Use(h.AuthMiddleware)
	h.Basic.Register(api)
}
