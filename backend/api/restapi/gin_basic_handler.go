package restapi

import (
	"github.com/gin-gonic/gin"
)

type BasicHandler struct {
}

func NewBasicHandler() *BasicHandler {
	return &BasicHandler{}
}

func (h *BasicHandler) Register(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
}
