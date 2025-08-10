package restapi

import (
	"net/http"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type BasicHandler struct {
	svgService service.SvgService
}

func NewBasicHandler(ss service.SvgService) *BasicHandler {
	return &BasicHandler{svgService: ss}
}

func (h *BasicHandler) Register(router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	router.GET("/user-svgs", h.getUserSvgs)
}

func (h *BasicHandler) getUserSvgs(c *gin.Context) {
	svgs, err := h.svgService.GetUserSvgs(c)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, svgs)
}
