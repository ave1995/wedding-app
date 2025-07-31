package restapi

import (
	"net/http"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) *userHandler {
	return &userHandler{userService: us}
}

func (h *userHandler) Register(router *gin.RouterGroup) {
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
//	@Success		201		{object}	domain.User
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/register [post]
func (h *userHandler) registerUser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithBadRequest(c, err, "invalid request body")
		return
	}

	user, err := h.userService.RegisterUser(c, req.Username, req.Email, req.Password)
	if err != nil {
		respondWithInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}
