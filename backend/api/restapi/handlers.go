package restapi

import "github.com/gin-gonic/gin"

type Handlers struct {
	User  *userHandler
	Basic *basicHandler
}

func (h *Handlers) RegisterAll(router *gin.Engine) {
	auth := router.Group("/auth")
	h.User.Register(auth)

	api := router.Group("/api")
	h.Basic.Register(api)
}
