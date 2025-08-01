package restapi

import (
	"github.com/gin-gonic/gin"
)

type basicHandler struct {
}

func NewBasicHandler() *basicHandler {
	return &basicHandler{}
}

func (h *basicHandler) Register(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
