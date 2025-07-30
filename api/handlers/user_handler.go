package handlers

import (
	"net/http"
	"wedding-app/handlers/responses"
	"wedding-app/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(us services.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterUserHandler godoc
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
func (h *UserHandler) RegisterUserHandler(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.RespondWithBadRequest(c, err, "invalid request body")
		return
	}

	user, err := h.userService.RegisterUser(c, req.Username, req.Email, req.Password)
	if err != nil {
		responses.RespondWithInternalError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}
