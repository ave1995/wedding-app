package restapi

import "github.com/gin-gonic/gin"

type Handlers struct {
	User *userHandler
}

func (h *Handlers) RegisterAll(router *gin.Engine) {
	auth := router.Group("/auth")
	h.User.Register(auth)
}
