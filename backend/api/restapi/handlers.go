package restapi

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handlers struct {
	User  *userHandler
	Basic *basicHandler
}

func (h *Handlers) RegisterAll(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	h.User.Register(auth)

	api := router.Group("/api")
	h.Basic.Register(api)
}
