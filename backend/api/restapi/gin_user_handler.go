package restapi

import (
	"net/http"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

func (h *UserHandler) Register(router *gin.RouterGroup) {
	router.POST("/register", h.registerUser)
}

// registerUser godoc
//
//	@Summary		Register a new user
//	@Description	Create a user with username, email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		RegisterRequest	true	"User info"
//	@Success		201		{object}	model.User
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/register [post]
func (h *UserHandler) registerUser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewInvalidRequestPayloadAPIError(err))
		return
	}

	user, err := h.userService.RegisterUser(c, req.Username, req.Email, req.Password)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusCreated, user)
}
