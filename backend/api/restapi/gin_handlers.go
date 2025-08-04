package restapi

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinHandlers struct {
	User  *UserHandler
	Basic *BasicHandler
}

func (h *GinHandlers) RegisterAll(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	h.User.Register(auth)

	api := router.Group("/api")
	h.Basic.Register(api)
}
