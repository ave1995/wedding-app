package restapi

import (
	"log/slog"
	"net/http"
	"time"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
	logger      *slog.Logger
}

func NewUserHandler(us service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{userService: us, logger: logger}
}

func (h *UserHandler) Register(router *gin.RouterGroup) {
	router.POST("/register", h.registerUser)
	router.POST("/login", h.loginUser)
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

// registerUser godoc
//
//	@Summary		Login a existing user
//	@Description	Login with with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		LoginRequest	true	"User info"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/auth/login [post]
func (h *UserHandler) loginUser(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(NewInvalidRequestPayloadAPIError(err))
		return
	}

	accessToken, err := h.userService.LoginUser(c, req.Email, req.Password)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	h.logger.Info("Login successful",
		slog.String("token", accessToken.Token),
		slog.Time("duration", accessToken.ExpiresAt))

	c.SetCookie(
		"access_token",
		accessToken.Token,
		int(time.Until(accessToken.ExpiresAt).Seconds()), // cookie lifetime
		"/",  // path
		"",   // domain ("" works for localhost)
		true, // secure (HTTPS only)
		true, // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}
